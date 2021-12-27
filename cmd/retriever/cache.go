package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
	"github.com/patrickmn/go-cache"
)

func getSnapshotItemsCached(ctx context.Context, goCache *cache.Cache, bucketName string) []util.SnapshotItem {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := util.NewSnapshotStorage(client, bucketName)
	items := []util.SnapshotItem{}

	itemsFound, found := goCache.Get("items")
	if found {
		items = itemsFound.([]util.SnapshotItem)
	}

	items = snapshotStorage.GetSnapshotItems(ctx)
	goCache.Set("items", items, cache.DefaultExpiration)

	return items
}

func getNewestSnapshotsPublicURL(
	ctx context.Context,
	goCache *cache.Cache,
	bucketName string,
	network util.NetworkType,
) (string, error) {
	items := getSnapshotItemsCached(ctx, goCache, bucketName)

	itemsNetwork := []util.SnapshotItem{}

	for _, item := range items {
		if item.Network == network {
			itemsNetwork = append(itemsNetwork, item)
		}
	}

	if len(itemsNetwork) < 1 {
		return "", fmt.Errorf("Snapshot item from %s network not found", network)
	}

	return itemsNetwork[0].Link, nil
}
