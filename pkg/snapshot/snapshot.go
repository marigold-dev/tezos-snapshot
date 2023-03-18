package snapshot

import (
	"math"
	"time"
)

type BlockHeaderResponse struct {
	Level     int    `json:"level"`
	Proto     int    `json:"proto"`
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
	ChainID   string `json:"chain_id"`
	Signature string `json:"signature"`
}

type SnapshotItem struct {
	FileName       string       `json:"file_name"`
	Chain          string       `json:"chain"`
	BlockTimestamp string       `json:"block_timestamp"`
	BlockHash      string       `json:"block_hash"`
	BlockHeight    string       `json:"block_height"`
	URL            string       `json:"url"`
	Filesize       string       `json:"filesize"`
	SHA256         string       `json:"sha256"`
	ArtifactType   string       `json:"artifact_type"`
	FilesizeBytes  int64        `json:"filesize_bytes"`
	Date           time.Time    `json:"date"`
	NetworkType    NetworkType  `json:"network_type"`
	SnapshotType   SnapshotType `json:"snapshot_type"`
	TezosVersion   TezosVersion `json:"tezos_version"`
}

type TezosVersion struct {
	Implementation string     `json:"implementation"`
	Version        Version    `json:"version"`
	CommitInfo     CommitInfo `json:"commit_info"`
}

type Version struct {
	Major          int    `json:"major"`
	Minor          int    `json:"minor"`
	AdditionalInfo string `json:"additional_info"`
}

type CommitInfo struct {
	CommitHash string `json:"commit_hash"`
	CommitDate string `json:"commit_date"`
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

func NetworkProtocolPriority(chain string) int {
	if chain == "ITHACA" {
		return 0
	}

	if chain == "MAIN" {
		return math.MaxInt
	}

	network := chain
	network_char := network[0]
	return int(network_char)
}
