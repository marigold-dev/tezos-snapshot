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

func (s *SnapshotHeader) SanitizeChainName() string {
	chainName := s.ChaiName
	if strings.Contains(chainName, "_") {
		chainName = strings.Split(chainName, "_")[1]
	}

	if strings.Contains(chainName, "-") {
		chainName = strings.Split(chainName, "-")[0]
	}

	chainName = strings.ToLower(chainName)

	if chainName == "ithacanet" {
		chainName = "ghostnet"
	}

	return chainName
}
