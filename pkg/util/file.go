package util

import (
	"log"
	"strings"
	"time"

	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
)

type File struct {
	Name string
}

func (f *File) YearMonth() int {
	date := f.Date()
	return ((int(date.Month()) * 12) + int(date.Year()))
}

func (f *File) Date() time.Time {
	paths := strings.Split(f.Name, "/")

	if len(paths) <= 0 {
		log.Fatalf("Invalid file name %q. \n", f.Name)
	}

	folderName := paths[0]
	date, err := time.Parse("2006.01.02", folderName)
	if err != nil {
		log.Fatalf("Invalid file name %q. \n", f.Name)
	}
	return date
}

func (f *File) NetworkProtocol() snapshot.NetworkProtocolType {
	return snapshot.NetworkProtocolType(strings.Split(strings.Split(f.Name, "-")[0], "_")[1])
}

func (f *File) SnapshotType() snapshot.SnapshotType {
	snapshotType := snapshot.SnapshotType(snapshot.FULL)

	if strings.Contains(f.Name, "rolling") {
		snapshotType = snapshot.SnapshotType(snapshot.ROLLING)
	}

	return snapshotType
}
