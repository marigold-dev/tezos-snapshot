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
	endpoint := os.Getenv("ENDPOINT")

	if bucketName == "" {
		log.Fatalln("The BUCKET_NAME environment variable is empty.")
	}

	if endpoint == "" {
		log.Fatalln("The ENDPOINT environment variable is empty.")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := util.NewSnapshotStorage(client, bucketName)

	createSnapshot(endpoint, isRollingSnapshot)

	snapshotfileName := getSnapshotNames(isRollingSnapshot)

	snapshotStorage.EphemeralUpload(ctx, snapshotfileName)

	snapshotStorage.DeleteOldSnapshots(ctx, maxDays)
}
