package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
)

func main() {
	ctx := context.Background()
	maxDays := util.GetEnvInt("MAX_DAYS", 7)
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatalln("The BUCKET_NAME environment variable is empty.")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := util.NewSnapshotStorage(client, bucketName)

	createSnapshot(false)
	createSnapshot(true)

	snapshotfileNameFull, snapshotfileNamesRolling := getSnapshotsNames()

	snapshotStorage.EphemeralUpload(ctx, snapshotfileNameFull, false)
	snapshotStorage.EphemeralUpload(ctx, snapshotfileNamesRolling, true)

	snapshotStorage.DeleteOldSnapshots(ctx, maxDays)
}
