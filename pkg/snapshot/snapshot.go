package snapshot

import "time"

type SnapshotItem struct {
	FileName        string
	Network         NetworkType
	NetworkProtocol string
	Date            time.Time
	SnapshotType    SnapshotType
	Blockhash       string
	Blocklevel      string
	PublicURL       string
	Size            int64
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
