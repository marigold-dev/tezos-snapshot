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
	Filename       string          `json:"filename"`
	ChainName      string          `json:"chain_name"`
	BlockTimestamp string          `json:"block_timestamp"`
	BlockHash      string          `json:"block_hash"`
	BlockHeight    int             `json:"block_height"`
	URL            string          `json:"url"`
	Filesize       string          `json:"filesize"`
	SHA256         string          `json:"sha256"`
	ArtifactType   ArtifactType    `json:"artifact_type"`
	HistoryMode    HistoryModeType `json:"history_mode"`
	FilesizeBytes  int64           `json:"filesize_bytes"`
	Date           time.Time       `json:"date"`
	TezosVersion   TezosVersion    `json:"tezos_version"`
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

type ArtifactType string
type HistoryModeType string

const (
	SNAPSHOT ArtifactType = "tezos-snapshot"
	TARBALL  ArtifactType = "tarball"
)
const (
	ROLLING HistoryModeType = "rolling"
	FULL    HistoryModeType = "full"
	ARCHIVE HistoryModeType = "archive"
)

// NetworkProtocolPriority it's a way to sort like that:
// 1. Mainnet
// 2. Limannet,
// 3. Mumbainet
// 4. Etc..
// 5. Ithacanet/Ghostnet
func NetworkProtocolPriority(chain string) int {
	// Ithacanet/Ghostnet, then will be the last on the list
	if chain == "ithaca" {
		return 0
	}

	// Mainnet then will be the first on the list
	if chain == "main" {
		return math.MaxInt
	}

	// Others protocol by protocol number
	network := chain
	network_char := network[0]
	return int(network_char)
}
