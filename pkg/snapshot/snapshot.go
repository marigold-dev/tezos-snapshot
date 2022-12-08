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
	MAIN         NetworkProtocolType = "MAINNET"
	HANGZHOU     NetworkProtocolType = "HANGZHOUNET"
	ITHACA       NetworkProtocolType = "ITHACANET"
	JAKARTA      NetworkProtocolType = "JAKARTANET"
	KATHMANDUNET NetworkProtocolType = "KATHMANDUNET"
	LIMANET      NetworkProtocolType = "LIMANET"
)

func NetworkProtocolPriority(networkProtocol NetworkProtocolType) int {
	if networkProtocol == ITHACA {
		return 0
	}

	if networkProtocol == MAIN {
		return math.MaxInt
	}

	network := string(networkProtocol)
	network_char := network[0]
	return int(network_char)
}
