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
	"github.com/samber/lo"
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
	log.Printf("Opening snapshot file %q.", fileName)
	snapshotFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer snapshotFile.Close()

	log.Printf("Uploading %q snapshot to Google Clound Storage.", fileName)
	err = s.uploadSnapshot(ctx, snapshotFile)
	if err != nil {
		log.Fatalf("%v \n", err)
	}

	log.Printf("Deleting snapshot file %q.", fileName)
	err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *SnapshotStorage) GetTodaySnapshotsItems(ctx context.Context) []snapshot.SnapshotItem {
	items := s.GetSnapshotItems(ctx)
	todayItems := lo.Filter(items, func(item snapshot.SnapshotItem, _ int) bool {
		now := time.Now()
		return (item.Date.YearDay() == now.YearDay() && item.Date.Year() == now.Year())
	})
	return todayItems
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
		if items[i].Date == items[j].Date {
			networkIsPriority :=
				snapshot.NetworkProtocolPriority(items[i].NetworkProtocol) >
					snapshot.NetworkProtocolPriority(items[j].NetworkProtocol)
			return networkIsPriority
		}
		dateIsGreater := items[i].Date.After(items[j].Date)
		return dateIsGreater
	})

	return items
}

func (s *SnapshotStorage) DeleteExpiredSnapshots(ctx context.Context, maxDays int, maxMonths int) {
	log.Println("Deleting expired snapshots in the Google Cloud Storage.")

	it := s.client.Bucket(s.bucketName).Objects(ctx, &storage.Query{})

	files := []File{}
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("listBucket: unable to list bucket %q: %v \n", s.bucketName, err)
		}

		file := File{
			Name: obj.Name,
		}

		files = append(files, file)
	}

	now := time.Now()
	filesToDelete := filterFilesToDelete(maxDays, maxMonths, files, now)

	lo.ForEach(filesToDelete, func(file File, _ int) {
		log.Printf("Deleting %q object. \n", file.Name)
		objHandler := s.client.Bucket(s.bucketName).Object(file.Name)
		err := objHandler.Delete(ctx)
		if err != nil {
			log.Printf("%v \n", err)
		}
		log.Printf("%q object deleted. \n", file.Name)
	})
}

func (s *SnapshotStorage) uploadSnapshot(ctx context.Context, file *os.File) error {
	hasher := sha256.New()
	currentTime := time.Now()
	currentDate := currentTime.Format("2006.01.02")

	log.Printf("Current Date is %q.\n", currentDate)

	fileName := currentDate + "/" + file.Name()

	objectHandler := s.client.Bucket(s.bucketName).Object(fileName)

	objWriter := objectHandler.NewWriter(ctx)

	writer := io.MultiWriter(objWriter, hasher)

	if _, err := io.Copy(writer, file); err != nil {
		log.Printf("Error Write Copy")
		return err
	}

	if err := objWriter.Close(); err != nil {
		log.Printf("Error Write Close")
		return err
	}
	log.Printf("Blob %q uploaded.\n", file.Name())

	// Make this file public
	acl := objectHandler.ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}
	log.Printf("Blob %q is public now.\n", file.Name())

	// Add Checksum Metadata
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"SHA256CHECKSUM": fmt.Sprintf("%x", hasher.Sum(nil)),
		},
	}

	if _, err := objectHandler.Update(ctx, objectAttrsToUpdate); err != nil {
		log.Printf("Error Update SHA256 checksum metadata")
		return err
	}

	return nil
}

func filterFilesToDelete(maxDays int, maxMonths int, files []File, now time.Time) []File {
	fmt.Printf("Current Date is %q.\n", now.Format("2006.01.02"))
	actualMounth := ((int(now.Month()) * 12) + int(now.Year()))

	filesByYearMonthLookUp := lo.GroupBy(files, func(file File) int {
		return file.YearMonth()
	})

	filesByProtocolPriorityLookUp := lo.GroupBy(files, func(file File) snapshot.NetworkProtocolType {
		return file.NetworkProtocol()
	})

	checkFileMustBeDeleted := func(file File, _ int) bool {
		log.Printf("Check if is needed delete %q object. \n", file.Name)

		// Files where its month it's not more than (maxMonth) months ago and
		// it's not the first snapshot from its month.
		if (file.YearMonth() - actualMounth) < maxMonths {
			fmt.Printf("File YearMonth %d", file.YearMonth())
			filesYearMonth := filesByYearMonthLookUp[file.YearMonth()]
			filesYearMonthSameProtocolAndType := lo.Filter(filesYearMonth, func(f File, _ int) bool {
				return file.NetworkProtocol() == f.NetworkProtocol() && file.SnapshotType() == f.SnapshotType()
			})

			fmt.Printf("here2: %v", filesYearMonthSameProtocolAndType)
			firstFileWithThisMonth := lo.MinBy(filesYearMonthSameProtocolAndType, func(item File, min File) bool {
				return item.Date().Before(min.Date())
			})

			fmt.Printf("here: %v", firstFileWithThisMonth)

			if file == firstFileWithThisMonth {
				return false
			}
		}

		// Files where is not the first protocols file
		filesProtocolPriority := filesByProtocolPriorityLookUp[file.NetworkProtocol()]
		filesProtocolPrioritySameProtocolAndType := lo.Filter(filesProtocolPriority, func(f File, _ int) bool {
			return file.NetworkProtocol() == f.NetworkProtocol() && file.SnapshotType() == f.SnapshotType()
		})
		firstFileWithThisProtocol := lo.MaxBy(filesProtocolPrioritySameProtocolAndType, func(item File, max File) bool {
			return item.Date().After(max.Date())
		})
		if file == firstFileWithThisProtocol {
			return false
		}

		// Files where is more than (maxDays) days ago

		diffDays := int(now.Sub(file.Date()).Hours() / 24)
		log.Printf("Date folder diffDays is %d. \n", diffDays)

		if diffDays <= maxDays {
			return false
		}

		log.Printf("Delete %q object. \n", file.Name)
		return true
	}

	return lo.Filter(files, checkFileMustBeDeleted)
}

func isFile(file *storage.ObjectAttrs) (bool, string, string) {
	splitedBySlash := strings.Split(file.Name, "/")
	if len(splitedBySlash) < 2 {
		return false, "", ""
	}

	return (len(splitedBySlash) == 2 && (splitedBySlash[1] != "")), splitedBySlash[0], splitedBySlash[1]
}
