package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
)

type SnapshotExec struct {
	snapshotsPath    string
	octezNodeBinPath string
	tezosPath        string
}

func NewSnapshotExec(snapshotsPath, octezNodePath, tezosPath string) *SnapshotExec {
	return &SnapshotExec{snapshotsPath, octezNodePath, tezosPath}
}

func (s *SnapshotExec) CreateSnapshot(historyMode snapshot.HistoryModeType) {
	log.Println("Creating snapshot.")
	script := "mkdir -p " + s.snapshotsPath + " && cd " + s.snapshotsPath + " && " + s.octezNodeBinPath + " snapshot export --block head~10 --data-dir " + s.tezosPath + "/data"

	if historyMode == snapshot.ROLLING {
		script = script + " --rolling"
	}

	_, _ = s.execScript(script)
}

func (s *SnapshotExec) GetSnapshotName(historyMode snapshot.HistoryModeType) (string, error) {
	log.Println("Getting snapshot names.")
	script := "mkdir -p " + s.snapshotsPath + " && cd " + s.snapshotsPath + " && /bin/ls -1a"
	stdout, _ := s.execScript(script)

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

func (s *SnapshotExec) GetSnapshotHeaderOutput(filepath string) string {
	log.Printf("Getting snapshot header output for file: %q. \n", filepath)
	script := s.octezNodeBinPath + " snapshot info --json " + s.snapshotsPath + "/" + filepath
	stdout, _ := s.execScript(script)
	log.Printf("Snapshot header output: %q. \n", stdout.String())
	return stdout.String()
}

func (s *SnapshotExec) DeleteLocalSnapshots() {
	log.Println("Deleting local snapshots.")
	script := "rm -rf " + s.snapshotsPath + "/*"
	_, _ = s.execScript(script)
}

func (s *SnapshotExec) execScript(script string) (bytes.Buffer, bytes.Buffer) {
	log.Printf("Executing script: %q. \n", script)
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
