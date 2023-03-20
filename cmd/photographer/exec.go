package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
)

func createSnapshot(snapshotType snapshot.SnapshotType) {
	bin := "/usr/local/bin/octez-node"

	args := []string{"snapshot", "export", "--block", "head~30", "--data-dir", "/var/run/tezos/node/data"}

	if snapshotType == snapshot.ROLLING {
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

func getSnapshotNames(snapshotType snapshot.SnapshotType) (string, error) {
	log.Println("Getting snapshot names.")
	var errBuf, outBuf bytes.Buffer
	cmd := exec.Command("/bin/ls", "-1a")
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotfilenames := strings.Split(outBuf.String(), "\n")
	log.Printf("All files found: %v \n", snapshotfilenames)

	extension := "full"

	if snapshotType == snapshot.ROLLING {
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
