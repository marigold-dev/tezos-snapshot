package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
)

func createSnapshot(historyMode snapshot.HistoryModeType) {
	bin := "/usr/local/bin/octez-node"

	args := []string{"sh", "-c", "mkdir -p /var/run/tezos/snapshots && cd /var/run/tezos/snapshots && snapshot export --block head~30 --data-dir /var/run/tezos/node/data"}

	if historyMode == snapshot.ROLLING {
		args = append(args, "--rolling")
	}

	cmd := exec.Command(bin, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Start()
	if err != nil {
		log.Fatalf("%v \n", err)
	}

	log.Println("snapshot export stdout:")
	log.Println(stdout.String())
}

func getSnapshotNames(historyMode snapshot.HistoryModeType) (string, error) {
	log.Println("Getting snapshot names.")
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "mkdir -p /var/run/tezos/snapshots && cd /var/run/tezos/snapshots && /bin/ls -1a")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotfilenames := strings.Split(stdout.String(), "\n")
	log.Printf("All files found: %v \n", snapshotfilenames)

	extension := "full"

	if historyMode == snapshot.ROLLING {
		extension = "rolling"
	}

	for _, filename := range snapshotfilenames {
		if strings.Contains(filename, extension) {
			log.Printf("Snapshot file found is: %q. \n", filename)
			return filename, nil
		}
	}

	return "", fmt.Errorf("Snapshot file not found.")
}
