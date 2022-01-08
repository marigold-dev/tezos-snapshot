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
	bucketName := os.Getenv("BUCKET_NAME")
	maxDays := util.GetEnvInt("MAX_DAYS", 7)
	isRollingSnapshot := util.GetEnvBool("ROLLING", false)

	if bucketName == "" {
		log.Fatalln("The BUCKET_NAME environment variable is empty.")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := util.NewSnapshotStorage(client, bucketName)

	createSnapshot(isRollingSnapshot)

	snapshotfileNameFull, snapshotfileNamesRolling := getSnapshotsNames()

	if isRollingSnapshot {
		snapshotStorage.EphemeralUpload(ctx, snapshotfileNamesRolling, isRollingSnapshot)
	} else {
		snapshotStorage.EphemeralUpload(ctx, snapshotfileNameFull, isRollingSnapshot)
	}

	snapshotStorage.DeleteOldSnapshots(ctx, maxDays)
}
