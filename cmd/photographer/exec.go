package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
)

func createSnapshot(rolling bool) {
	bin := "/usr/local/bin/tezos-node"

	args := []string{"snapshot", "export", "--block", "head~30", "--data-dir", "/var/run/tezos/node/data"}

	if rolling {
		args = append(args, "--rolling")
	}

	cmd := exec.Command(bin, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("%v \n", err)
	}
}

func getSnapshotNames(isRolling bool) string {
	log.Println("Getting snapshot names.")
	var errBuf, outBuf bytes.Buffer
	cmd := exec.Command("/bin/ls", "-1a")
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotfileNames := strings.Split(outBuf.String(), "\n")

	log.Print(outBuf.String())
	log.Printf("len: %d \n", len(snapshotfileNames))

	var rolling, full string

	for _, fileName := range snapshotfileNames {
		var fileNameLower = strings.ToLower(fileName)
		if strings.Contains(fileNameLower, "tezos") {
			rolling = fileNameLower
		}
	}

	if isRolling {
		log.Printf("Rolling snapshot file is: %q. \n", rolling)
		return rolling
	}

	log.Printf("Full snapshot file is: %q. \n", full)
	return full
}

func execute(ctx context.Context, snapshotStorage *util.SnapshotStorage, rolling bool, network snapshot.NetworkProtocolType) {
	todayItems := snapshotStorage.GetTodaySnapshotsItems(ctx)

	snapshotType := snapshot.FULL

	if rolling {
		snapshotType = snapshot.ROLLING
	}

	alreadyExist := util.Some(todayItems, func(item snapshot.SnapshotItem) bool {
		return item.NetworkProtocol == network && item.SnapshotType == snapshotType
	})

	if alreadyExist {
		return
	}

	createSnapshot(rolling)
	snapshotfileName := getSnapshotNames(rolling)
	snapshotStorage.EphemeralUpload(ctx, snapshotfileName)
}
