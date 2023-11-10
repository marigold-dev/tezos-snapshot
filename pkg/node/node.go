package node

import (
	"io"
	"log"
	"net/http"
	"time"
)

func CheckNodesAreReady() {
	for {
		r, err := http.Get("http://localhost:8732/version")
		if err != nil && r.StatusCode != 200 {
			log.Println("The node is not running. Waiting 5 minutes...")
			time.Sleep(5 * time.Minute)
		}
		defer r.Body.Close()
		if r.StatusCode == 200 {
			break
		}
	}
}

func GetTezosVersion() string {
	reqVersion, err := http.Get("http://localhost:8732/version")
	if err != nil {
		log.Fatalf("Unable to get node version. %v \n", err)
	}
	defer reqVersion.Body.Close()
	version, err := io.ReadAll(reqVersion.Body)
	if err != nil {
		log.Fatalf("Unable to read node version. %v \n", err)
	}

	return string(version)
}
