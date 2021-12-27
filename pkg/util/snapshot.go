package util

import "time"

type SnapshotItem struct {
	FileName     string
	Network      NetworkType
	Date         time.Time
	SnapshotType SnapshotType
	Blockhash    string
	Blocklevel   string
	Link         string
}

type SnapshotType string
type NetworkType string

const (
	ROLLING SnapshotType = "ROLLING"
	FULL    SnapshotType = "FULL"
)
const (
	MAINNET NetworkType = "MAINNET"
	TESTNET NetworkType = "TESTNET"
)
