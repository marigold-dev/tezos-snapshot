package snapshot

import (
	"encoding/json"
	"strings"
)

type SnapshotHeader struct {
	Version   int    `json:"version"`
	ChaiName  string `json:"chain_name"`
	Mode      string `json:"mode"`
	BlockHash string `json:"block_hash"`
	Level     int    `json:"level"`
	Timestamp string `json:"timestamp"`
}

type WrapperSnapshotHeader struct {
	SnapshotHeader SnapshotHeader `json:"snapshot_header"`
}

func SnapshotHeaderFromJson(snapshotHeaderOutput string) (*SnapshotHeader, error) {
	var snapshotHeader WrapperSnapshotHeader
	err := json.Unmarshal([]byte(snapshotHeaderOutput), &snapshotHeader)
	if err != nil {
		return nil, err
	}

	return &snapshotHeader.SnapshotHeader, nil
}

// Example: TEZOS_ITHACANET_2022-01-25T15:00:00Z to ghostnet
// Example: TEZOS_MAINNETrolling to mainnet
func (s *SnapshotHeader) SanitizeChainame() string {
	parts := strings.Split(s.ChaiName, "_")
	chainName := strings.ToLower(parts[1])
	if chainName == "ithacanet" {
		chainName = "ghostnet"
	}
	return chainName
}
