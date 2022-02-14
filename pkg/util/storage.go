package util

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"google.golang.org/api/iterator"
)

type SnapshotStorage struct {
	client     *storage.Client
	bucketName string
}

func NewSnapshotStorage(client *storage.Client, bucketName string) *SnapshotStorage {
	return &SnapshotStorage{client: client, bucketName: bucketName}
}

func (s *SnapshotStorage) EphemeralUpload(ctx context.Context, fileName string) {
	fmt.Printf("Opening snapshot file %q.", fileName)
	snapshotFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer snapshotFile.Close()

	fmt.Printf("Uploading %q snapshot to Google Clound Storage.", fileName)
	err = s.uploadSnapshot(ctx, snapshotFile)
	if err != nil {
		log.Fatalf("%v \n", err)
	}

	fmt.Printf("Deleting snapshot file %q.", fileName)
	err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SnapshotStorage) GetSnapshotItems(ctx context.Context) []snapshot.SnapshotItem {

	items := []snapshot.SnapshotItem{}
	it := s.client.Bucket(s.bucketName).Objects(ctx, &storage.Query{})

	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("listBucket: unable to list bucket %q: %v \n", s.bucketName, err)
		}

		isFile, folderName, fileName := isFile(obj)

		layout := "2006.01.02"
		date, err := time.Parse(layout, folderName)
		if err != nil {
			log.Fatalf("unable to parse date. %v \n", err)
		}

		if !isFile {
			continue
		}

		network := snapshot.NetworkType(snapshot.TESTNET)
		networkProtocol := snapshot.NetworkProtocolType(strings.Split(strings.Split(fileName, "-")[0], "_")[1])
		size := obj.Size

		if strings.Contains(obj.Name, "MAINNET") {
			network = snapshot.NetworkType(snapshot.MAINNET)
		}
		snapshotType := snapshot.SnapshotType(snapshot.FULL)

		if strings.Contains(obj.Name, "rolling") {
			snapshotType = snapshot.SnapshotType(snapshot.ROLLING)
		}

		splitedByHyphen := strings.Split(obj.Name, "-")

		blocklevel := splitedByHyphen[len(splitedByHyphen)-1]
		blockhash := splitedByHyphen[len(splitedByHyphen)-2]

		checksum := obj.Metadata["SHA256CHECKSUM"]

		item := snapshot.SnapshotItem{
			FileName:        fileName,
			Network:         network,
			Size:            size,
			NetworkProtocol: networkProtocol,
			Date:            date,
			SnapshotType:    snapshotType,
			PublicURL:       obj.MediaLink,
			Blockhash:       blockhash,
			Blocklevel:      blocklevel,
			SHA256Checksum:  checksum,
		}

		items = append(items, item)
	}

	// Order by date and network
	sort.Slice(items, func(i, j int) bool {
		dateIsGreater := items[i].Date.After(items[j].Date)
		networkIsPriority :=
			snapshot.NetworkProtocolPriority(items[i].NetworkProtocol) >
				snapshot.NetworkProtocolPriority(items[j].NetworkProtocol)

		return dateIsGreater && networkIsPriority
	})

	return items
}

func (s *SnapshotStorage) DeleteOldSnapshots(ctx context.Context, maxDays int) {
	fmt.Println("Deleting old snapshots in the Google Cloud Storage.")

	it := s.client.Bucket(s.bucketName).Objects(ctx, &storage.Query{})

	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("listBucket: unable to list bucket %q: %v \n", s.bucketName, err)
		}

		err = s.deleteFile(ctx, maxDays, obj)
		if err != nil {
			fmt.Printf("%v \n", err)
		}
	}
}

func (s *SnapshotStorage) uploadSnapshot(ctx context.Context, file *os.File) error {
	hasher := sha256.New()
	currentTime := time.Now()
	currentDate := currentTime.Format("2006.01.02")

	fmt.Printf("Current Date is %q.\n", currentDate)

	fileName := currentDate + "/" + file.Name()

	objectHandler := s.client.Bucket(s.bucketName).Object(fileName)

	objWriter := objectHandler.NewWriter(ctx)

	writer := io.MultiWriter(objWriter, hasher)

	if _, err := io.Copy(writer, file); err != nil {
		fmt.Printf("Error Write Copy")
		return err
	}

	if err := objWriter.Close(); err != nil {
		fmt.Printf("Error Write Close")
		return err
	}
	fmt.Printf("Blob %q uploaded.\n", file.Name())

	// Make this file public
	acl := objectHandler.ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	fmt.Printf("Blob %q is public now.\n", file.Name())

	// Add Checksum Metadata
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"SHA256CHECKSUM": fmt.Sprintf("%x", hasher.Sum(nil)),
		},
	}

	if _, err := objectHandler.Update(ctx, objectAttrsToUpdate); err != nil {
		fmt.Printf("Error Update SHA256 checksum metadata")
		return err
	}

	return nil
}

func (s *SnapshotStorage) deleteFile(ctx context.Context, maxDays int, obj *storage.ObjectAttrs) error {
	fmt.Printf("Check if is needed delete %q object. \n", obj.Name)

	paths := strings.Split(obj.Name, "/")

	if len(paths) <= 0 {
		return fmt.Errorf("Invalid file name %q. \n", obj.Name)
	}

	folderName := paths[0]
	fmt.Printf("Name folder is %q. \n", folderName)

	t, err := time.Parse("2006.01.02", folderName)
	if err != nil {
		return err
	}
	fmt.Printf("Date folder is %v. \n", t)

	diff := time.Now().Sub(t)
	fmt.Printf("Date folder diff is %d. \n", diff)

	diffDays := int(diff.Hours() / 24)
	fmt.Printf("Date folder diffDays is %d. \n", diffDays)

	if diffDays > maxDays {
		fmt.Printf("Deleting %q object. \n", obj.Name)

		objHandler := s.client.Bucket(s.bucketName).Object(obj.Name)
		err = objHandler.Delete(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%q object deleted. \n", obj.Name)
	}
	return nil
}

func isFile(file *storage.ObjectAttrs) (bool, string, string) {
	splitedBySlash := strings.Split(file.Name, "/")
	if len(splitedBySlash) < 2 {
		return false, "", ""
	}

	return (len(splitedBySlash) == 2 && (splitedBySlash[1] != "")), splitedBySlash[0], splitedBySlash[1]
}
