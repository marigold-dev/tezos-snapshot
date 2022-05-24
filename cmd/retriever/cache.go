package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
	"github.com/patrickmn/go-cache"
)

func getSnapshotItemsCached(ctx context.Context, goCache *cache.Cache, bucketName string) []snapshot.SnapshotItem {
	itemsFound, found := goCache.Get("items")
	if found {
		log.Println("Using items from cache...")
		return itemsFound.([]snapshot.SnapshotItem)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := util.NewSnapshotStorage(client, bucketName)
	items := snapshotStorage.GetSnapshotItems(ctx)
	goCache.Set("items", items, cache.DefaultExpiration)
	return items
}

func getNewestSnapshot(
	ctx context.Context,
	goCache *cache.Cache,
	bucketName string,
	network snapshot.NetworkType,
	snapshotType snapshot.SnapshotType,
	protocol snapshot.NetworkProtocolType,
) (*snapshot.SnapshotItem, error) {
	items := getSnapshotItemsCached(ctx, goCache, bucketName)

	for _, item := range items {
		if item.Network == network && item.SnapshotType == snapshotType && item.NetworkProtocol == protocol {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("Snapshot item from %s network not found", network)
}
