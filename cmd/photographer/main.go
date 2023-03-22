package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
	"github.com/samber/lo"
)

func main() {
	godotenv.Load()

	cron := os.Getenv("CRON_EXPRESSION")
	if cron == "" {
		task()
	} else {
		log.Println("Waiting for the snapshot job...")
		s := gocron.NewScheduler(time.UTC)
		s.Cron("0 0 * * *").Do(task)
		s.StartBlocking()
	}
}

func task() {
	log.Println("Starting the snapshot job...")
	start := time.Now()
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")
	maxDays := util.GetEnvInt("MAX_DAYS", 7)
	maxMonths := util.GetEnvInt("MAX_MONTHS", 6)
	network := strings.ToUpper(os.Getenv("NETWORK"))

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

	// Check if today the rolling snapshot already exists
	execute(ctx, snapshotStorage, snapshot.ROLLING, network)

	// Check if today the full snapshot already exists
	execute(ctx, snapshotStorage, snapshot.FULL, network)

	snapshotStorage.DeleteExpiredSnapshots(ctx, maxDays, maxMonths)

	log.Printf("Snapshot job took %s", time.Since(start))
}

func execute(ctx context.Context, snapshotStorage *util.SnapshotStorage, snapshotType snapshot.SnapshotType, chain string) {
	todayItems := snapshotStorage.GetTodaySnapshotsItems(ctx)

	alreadyExist := lo.SomeBy(todayItems, func(item snapshot.SnapshotItem) bool {
		return item.ChainName == chain && item.SnapshotType == snapshotType
	})

	if alreadyExist {
		log.Printf("Already exist a today snapshot with %s type. \n", chain)
		return
	}

	createSnapshot(snapshotType)
	snapshotfilename, err := getSnapshotNames(snapshotType)
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotStorage.EphemeralUpload(ctx, snapshotfilename)
}
