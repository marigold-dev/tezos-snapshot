package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type SnapshotStorage struct {
	client     *storage.Client
	bucketName string
}

func NewSnapshotStorage(client *storage.Client, bucketName string) *SnapshotStorage {
	return &SnapshotStorage{client: client, bucketName: bucketName}
}

func (s *SnapshotStorage) EphemeralUpload(ctx context.Context, fileName string, rolling bool) {
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

func (s *SnapshotStorage) GetSnapshotItems(ctx context.Context) []SnapshotItem {

	items := []SnapshotItem{}
	it := s.client.Bucket(s.bucketName).Objects(ctx, &storage.Query{})

	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("listBucket: unable to list bucket %q: %v \n", s.bucketName, err)
		}

		isFile, folderName := isNotFolder(obj)
		layout := "2006-01-02"
		date, err := time.Parse(layout, folderName)

		if !isFile {
			continue
		}

		network := NetworkType(MAINNET)

		if strings.Contains(obj.Name, "TESTNET") {
			network = NetworkType(TESTNET)
		}

		snapshotType := SnapshotType(FULL)

		if strings.Contains(obj.Name, "rolling") {
			snapshotType = SnapshotType(ROLLING)
		}

		splited := strings.Split(obj.Name, "-")

		blocklevel := splited[len(splited)-1]
		blockhash := splited[len(splited)-2]

		item := SnapshotItem{
			FileName:     obj.Name,
			Network:      network,
			Date:         date,
			SnapshotType: snapshotType,
			Link:         obj.MediaLink,
			Blockhash:    blockhash,
			Blocklevel:   blocklevel,
		}

		items = append(items, item)
	}

	// Order by date
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.After(items[j].Date)
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
	currentTime := time.Now()
	currentDate := currentTime.Format("2006.01.02")

	fmt.Printf("Current Date is %q.\n", currentDate)

	objectHandler := s.client.Bucket(s.bucketName).Object(currentDate + "/" + file.Name())
	writer := objectHandler.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		fmt.Printf("Error Write Copy")
		return err
	}
	if err := writer.Close(); err != nil {
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

func isNotFolder(file *storage.ObjectAttrs) (bool, string) {
	splittedPaths := strings.Split(file.Name, "/")
	if len(splittedPaths) > 0 {
		return false, ""
	}

	return (len(splittedPaths) == 2 && (splittedPaths[1] != "")), splittedPaths[1]
}
