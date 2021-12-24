package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()
	maxDays := getEnvInt("MAX_DAYS", 7)
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatalln("The BUCKET_NAME environment variable is empty.")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := NewSnapshotStorage(client, bucketName)
	snapshotExec := NewSnapshotExec()

	snapshotExec.CreateSnapshot(false)
	snapshotExec.CreateSnapshot(true)

	snapshotfileNameFull, snapshotfileNamesRolling := snapshotExec.GetSnapshotsNames()

	snapshotStorage.EphemeralUpload(ctx, snapshotfileNameFull, false)
	snapshotStorage.EphemeralUpload(ctx, snapshotfileNamesRolling, true)

	snapshotStorage.DeleteOldSnapshots(ctx, maxDays)
}
