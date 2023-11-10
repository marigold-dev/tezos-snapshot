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

func SnapshotHeaderFromJson(snapshotHeaderOutput string) (*SnapshotHeader, error) {
	var snapshotHeader SnapshotHeader
	err := json.Unmarshal([]byte(snapshotHeaderOutput), &snapshotHeader)
	if err != nil {
		return nil, err
	}

	return &snapshotHeader, nil
}

// Example: TEZOS_MAINNET_2021-01-01_00-00 to mainnet
func (s *SnapshotHeader) SanitizeChainame() string {
	parts := strings.Split(s.ChaiName, "_")
	chainName := strings.ToLower(parts[len(parts)-1])
	if chainName == "ithacanet" {
		chainName = "ghostnet"
	}
	return chainName
}
