package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")
	maxDays := util.GetEnvInt("MAX_DAYS", 7)
	network := snapshot.NetworkProtocolType(strings.ToUpper(os.Getenv("NETWORK")))

	if bucketName == "" {
		log.Fatalln("The BUCKET_NAME environment variable is empty.")
	}

	if network == "" {
		log.Fatalln("The NETWORK environment variable is empty.")
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := util.NewSnapshotStorage(client, bucketName)

	// Check if the today rolling snapshot already exists
	execute(ctx, snapshotStorage, true, network)

	// Check if the today full snapshot already exists
	execute(ctx, snapshotStorage, false, network)

	snapshotStorage.DeleteOldSnapshots(ctx, maxDays)

	log.Printf("Snapshot job took %s", time.Since(start))
}
