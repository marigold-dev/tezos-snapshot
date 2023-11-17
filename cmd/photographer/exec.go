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
	script := "mkdir -p /var/run/tezos/snapshots && cd /var/run/tezos/snapshots && /usr/local/bin/octez-node snapshot export --block head~30 --data-dir /var/run/tezos/node/data"

	if historyMode == snapshot.ROLLING {
		script = script + " --rolling"
	}

	_, _ = execScript(script)
}

func getSnapshotName(historyMode snapshot.HistoryModeType) (string, error) {
	log.Println("Getting snapshot names.")
	script := "mkdir -p /var/run/tezos/snapshots && cd /var/run/tezos/snapshots && /bin/ls -1a"
	stdout, _ := execScript(script)

	snapshotfilenames := strings.Split(stdout.String(), "\n")
	log.Printf("All files found: %v \n", snapshotfilenames)

	for _, filename := range snapshotfilenames {
		if strings.Contains(filename, string(historyMode)) {
			log.Printf("Snapshot file found is: %q. \n", filename)
			return filename, nil
		}
	}

	return "", fmt.Errorf("Snapshot file not found.")
}

func getSnapshotHeaderOutput(filepath string) string {
	log.Printf("Getting snapshot header output for file: %q. \n", filepath)
	script := "/usr/local/bin/octez-node snapshot info --json" + filepath
	stdout, _ := execScript(script)
	log.Printf("Snapshot header output: %q. \n", stdout.String())
	return stdout.String()
}

func execScript(script string) (bytes.Buffer, bytes.Buffer) {
	cmd := exec.Command("sh", "-c", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	if stdout.Len() > 0 {
		log.Printf("stdout: \n%s\n", stdout.String())
	}
	if stderr.Len() > 0 {
		log.Printf("stderr: \n%s\n", stderr.String())
	}

	return stdout, stderr
}
