package store

import (
	"log"
	"strconv"
	"strings"

	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
)

type FileInfo struct {
	Filename    string
	ChainName   string
	HistoryMode snapshot.HistoryModeType
	BlockHeight int
	BlockHash   string
}

func getInfoFromfilename(filename string) *FileInfo {
	chainName := strings.ToLower(strings.Split(strings.Split(filename, "-")[0], "_")[1])

	if chainName == "ithacanet" {
		chainName = "ghostnet"
	}

	historyMode := snapshot.HistoryModeType(snapshot.FULL)

	if strings.Contains(filename, "rolling") {
		historyMode = snapshot.HistoryModeType(snapshot.ROLLING)
	}

	splitedByHyphen := strings.Split(filename, "-")
	blockheight, err := strconv.Atoi(strings.Split(splitedByHyphen[len(splitedByHyphen)-1], ".")[0])
	if err != nil {
		log.Fatalf("Unable to parse blockheight. %v \n", err)
	}
	blockhash := splitedByHyphen[len(splitedByHyphen)-2]

	return &FileInfo{
		Filename:    filename,
		ChainName:   chainName,
		HistoryMode: historyMode,
		BlockHeight: blockheight,
		BlockHash:   blockhash,
	}
}
