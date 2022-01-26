package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func createSnapshot(endpoint string, rolling bool) {
	bin := "/usr/local/bin/tezos-node"

	blockHash := relativeBlockHash(30, endpoint)
	fmt.Printf("block hash found: %s", blockHash)

	args := []string{"snapshot", "export", "--block " + blockHash, "--data-dir", "/var/run/tezos/node/data"}

	if rolling {
		args = append(args, "--rolling")
	}

	var errBuf, outBuf bytes.Buffer
	cmd := exec.Command(bin, args...)
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
}

func getSnapshotNames(isRolling bool) string {
	fmt.Println("Getting snapshot names.")
	var errBuf, outBuf bytes.Buffer
	cmd := exec.Command("/bin/ls", "-1a")
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotfileNames := strings.Split(outBuf.String(), "\n")

	fmt.Print(outBuf.String())
	fmt.Printf("len: %d \n", len(snapshotfileNames))

	var rolling, full string

	for _, fileName := range snapshotfileNames {
		if strings.Contains(fileName, "rolling") {
			rolling = fileName
		}
		if strings.Contains(fileName, "full") {
			full = fileName
		}
	}

	fmt.Printf("Full snapshot file is: %q. \n", full)
	fmt.Printf("Rolling snapshot file is: %q. \n", rolling)

	if isRolling {
		return rolling
	}

	return full
}

func relativeBlockHash(relative int, endpoint string) string {
	regex, err := regexp.Compile("(\"|')(.*)(\"|')")
	if err != nil {
		log.Fatalf("%v \n", err)
	}

	bin := "/usr/local/bin/tezos-client"

	args := []string{
		"--endpoint", endpoint,
		"rpc", "get",
		fmt.Sprintf("%s%d%s", "/chains/main/blocks/head~", relative, "/hash"),
	}

	cmd := exec.Command(bin, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%v \n", err)
		return ""
	}
	strOut := string(output)
	fmt.Println("Output getting hash block:")
	fmt.Println(strOut)

	regexResult := regex.FindString(strOut)
	regexResultWithoutSimpleQuotes := strings.ReplaceAll(regexResult, "'", "")
	regexResultWithoutQuotes := strings.ReplaceAll(regexResultWithoutSimpleQuotes, "\"", "")

	return regexResultWithoutQuotes
}
