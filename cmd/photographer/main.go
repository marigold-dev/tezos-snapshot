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
	"github.com/samber/lo"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")
	maxDays := util.GetEnvInt("MAX_DAYS", 7)
	maxMonths := util.GetEnvInt("MAX_MONTHS", 6)
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

	snapshotStorage.DeleteExpiredSnapshots(ctx, maxDays, maxMonths)

	log.Printf("Snapshot job took %s", time.Since(start))
}

func execute(ctx context.Context, snapshotStorage *util.SnapshotStorage, rolling bool, network snapshot.NetworkProtocolType) {
	todayItems := snapshotStorage.GetTodaySnapshotsItems(ctx)

	snapshotType := snapshot.FULL

	if rolling {
		snapshotType = snapshot.ROLLING
	}

	alreadyExist := lo.SomeBy(todayItems, func(item snapshot.SnapshotItem) bool {
		return item.NetworkProtocol == network && item.SnapshotType == snapshotType
	})

	if alreadyExist {
		log.Printf("Already exist a today snapshot with %q type. \n", network)
		return
	}

	createSnapshot(rolling)
	snapshotfileName, err := getSnapshotNames(rolling)
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotStorage.EphemeralUpload(ctx, snapshotfileName)
}
