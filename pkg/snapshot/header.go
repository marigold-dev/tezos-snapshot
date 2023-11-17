package snapshot

import (
	"encoding/json"
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
