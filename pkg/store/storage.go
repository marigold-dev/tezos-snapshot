package store

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/marigold-dev/tezos-snapshot/pkg/node"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
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

func (s *SnapshotStorage) EphemeralUpload(ctx context.Context, filename, snapshotHeaderOutput string) {
	// Save the current working directory so we can revert back
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	// Change to the desired directory
	if err := os.Chdir("/var/run/tezos/snapshots/"); err != nil {
		log.Fatalf("Failed to change directory: %v", err)
	}

	log.Printf("Opening snapshot file %q.", filename)
	snapshotFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	defer snapshotFile.Close()

	log.Printf("Uploading %q snapshot to Google Clound Storage.", filename)
	err = s.uploadSnapshot(ctx, snapshotFile, snapshotHeaderOutput)
	if err != nil {
		log.Fatalf("%v \n", err)
	}

	log.Printf("Deleting snapshot file %q.", filename)
	err = os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Chdir(currentDir) // Ensure we change back to the original directory
}

func (s *SnapshotStorage) GetTodaySnapshotsItems(ctx context.Context) []snapshot.SnapshotItem {
	items := s.GetSnapshotItems(ctx)
	now := time.Now()
	todayItems := lo.Filter(items, func(item snapshot.SnapshotItem, _ int) bool {
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

		isFile, folderName, filename := cloudObjIsFile(obj)

		layout := "2006.01.02"
		date, err := time.Parse(layout, folderName)
		if err != nil {
			log.Fatalf("unable to parse date. %v \n", err)
		}

		if !isFile {
			continue
		}

		size := obj.Size

		checksum := obj.Metadata["SHA256CHECKSUM"]
		snapshotHeaderJson := obj.Metadata["SNAPSHOT_HEADER"]
		versionJson, versionExist := obj.Metadata["VERSION"]
		version := snapshot.TezosVersion{}
		version.Implementation = "octez"
		if versionExist {
			json.Unmarshal([]byte(versionJson), &version)
		} else {
			version.Version.Major = 7
		}

		var snapshotHeader *snapshot.SnapshotHeader
		if snapshotHeaderJson == "" {
			// Handle old snapshots before the snapshot header was added
			filenameInfo := getInfoFromfilename(filename)
			snapshotVersion := 0
			if version.Version.Major >= 16 && version.Version.Major <= 17 {
				snapshotVersion = 5
			}
			if version.Version.Major == 18 && version.Version.Minor == 0 {
				snapshotVersion = 6
			}

			snapshotHeader = &snapshot.SnapshotHeader{
				Version:   snapshotVersion,
				ChaiName:  filenameInfo.ChainName,
				Mode:      string(filenameInfo.HistoryMode),
				BlockHash: filenameInfo.BlockHash,
				Level:     filenameInfo.BlockHeight,
				Timestamp: obj.Metadata["TIMESTAMP"],
			}
		} else {
			var err error
			snapshotHeader, err = snapshot.SnapshotHeaderFromJson(snapshotHeaderJson)
			if err != nil {
				log.Fatalf("Unable to parse snapshot header. %v \n", err)
			}
		}

		item := snapshot.SnapshotItem{
			Filename:        filename,
			Filesize:        util.FileSize(size),
			FilesizeBytes:   size,
			ChainName:       RemoveTezosPrefix(snapshotHeader.ChaiName),
			Date:            date,
			BlockTimestamp:  snapshotHeader.Timestamp,
			URL:             obj.MediaLink,
			BlockHash:       snapshotHeader.BlockHash,
			BlockHeight:     snapshotHeader.Level,
			SHA256:          checksum,
			TezosVersion:    version,
			ArtifactType:    snapshot.SNAPSHOT,
			HistoryMode:     snapshot.HistoryModeType(snapshotHeader.Mode),
			SnapshotVersion: snapshotHeader.Version,
		}

		items = append(items, item)
	}

	// Order by date and network
	sort.Slice(items, func(i, j int) bool {
		if items[i].Date == items[j].Date {
			networkIsPriority := items[i].NetworkProtocolPriority() > items[j].NetworkProtocolPriority()
			return networkIsPriority
		}
		dateIsGreater := items[i].Date.After(items[j].Date)
		return dateIsGreater
	})

	return items
}

func RemoveTezosPrefix(filename string) string {
	if strings.HasPrefix(filename, "TEZOS_") {
		return strings.ToLower(filename[6:])
	}
	return filename
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

func (s *SnapshotStorage) uploadSnapshot(ctx context.Context, file *os.File, snapshotHeaderOutput string) error {
	hasher := sha256.New()
	currentTime := time.Now()
	currentDate := currentTime.Format("2006.01.02")

	log.Printf("Current Date is %q.\n", currentDate)

	filename := currentDate + "/" + file.Name()

	objectHandler := s.client.Bucket(s.bucketName).Object(filename)

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

	// Wait nodes to be ready
	node.CheckNodesAreReady()

	// Request node version
	version := node.GetTezosVersion()

	// Getting Sha256 checksum
	sha256Checksum := fmt.Sprintf("%x", hasher.Sum(nil))

	// Add Checksum Metadata
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"SHA256CHECKSUM":  sha256Checksum,
			"VERSION":         version,
			"SNAPSHOT_HEADER": snapshotHeaderOutput,
		},
	}

	if _, err := objectHandler.Update(ctx, objectAttrsToUpdate); err != nil {
		log.Printf("Error Update SHA256 checksum metadata")
		return err
	}

	return nil
}

func filterFilesToDelete(maxDays int, maxMonths int, files []File, now time.Time) []File {
	log.Printf("Current Date is %q.\n", now.Format("2006.01.02"))
	actualMounth := ((int(now.Month()) * 12) + int(now.Year()))

	filesByYearMonthLookUp := lo.GroupBy(files, func(file File) int {
		return file.YearMonth()
	})

	filesByProtocolPriorityLookUp := lo.GroupBy(files, func(file File) string {
		return file.NetworkProtocol()
	})

	checkFileMustBeDeleted := func(file File, _ int) bool {
		log.Printf("Check if is needed delete %q object. \n", file.Name)

		// Files where its month it's not more than (maxMonth) months ago and
		// it's not the first snapshot from its month.
		if (file.YearMonth() - actualMounth) < maxMonths {
			log.Printf("File YearMonth %d", file.YearMonth())
			filesYearMonth := filesByYearMonthLookUp[file.YearMonth()]
			filesYearMonthSameProtocolAndType := lo.Filter(filesYearMonth, func(f File, _ int) bool {
				return file.NetworkProtocol() == f.NetworkProtocol() && file.HistoryMode() == f.HistoryMode()
			})

			firstFileWithThisMonth := lo.MinBy(filesYearMonthSameProtocolAndType, func(item File, min File) bool {
				return item.Date().Before(min.Date())
			})

			if file == firstFileWithThisMonth {
				return false
			}
		}

		// Files where is not the first protocols file
		filesProtocolPriority := filesByProtocolPriorityLookUp[file.NetworkProtocol()]
		filesProtocolPrioritySameProtocolAndType := lo.Filter(filesProtocolPriority, func(f File, _ int) bool {
			return file.NetworkProtocol() == f.NetworkProtocol() && file.HistoryMode() == f.HistoryMode()
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

func cloudObjIsFile(obj *storage.ObjectAttrs) (bool, string, string) {
	splitedBySlash := strings.Split(obj.Name, "/")
	if len(splitedBySlash) < 2 {
		return false, "", ""
	}

	return (len(splitedBySlash) == 2 && (splitedBySlash[1] != "")), splitedBySlash[0], splitedBySlash[1]
}
