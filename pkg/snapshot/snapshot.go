package snapshot

import (
	"math"
	"time"
)

type SnapshotItem struct {
	FileName        string
	Network         NetworkType
	NetworkProtocol NetworkProtocolType
	Date            time.Time
	SnapshotType    SnapshotType
	Blockhash       string
	Blocklevel      string
	PublicURL       string
	Size            int64
	SHA256Checksum  string
}

type SnapshotType string
type NetworkType string
type NetworkProtocolType string

const (
	ROLLING SnapshotType = "ROLLING"
	FULL    SnapshotType = "FULL"
)
const (
	MAINNET NetworkType = "MAINNET"
	TESTNET NetworkType = "TESTNET"
)

const (
	MAIN     NetworkProtocolType = "MAINNET"
	JAKARTA  NetworkProtocolType = "JAKARTA"
	HANGZHOU NetworkProtocolType = "HANGZHOUNET"
	ITHACA   NetworkProtocolType = "ITHACANET"
)

func NetworkProtocolPriority(networkProtocol NetworkProtocolType) int {
	switch networkProtocol {
	case MAIN:
		return 0
	case HANGZHOU:
		return 1
	case ITHACA:
		return 2
	}

	return math.MaxInt
}
