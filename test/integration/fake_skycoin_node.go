package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/simelo/rextporter/src/util/file"
	"github.com/simelo/rextporter/test/integration/testrand"
	log "github.com/sirupsen/logrus"
)

var jsonResponse = `
{
    "blockchain": {
        "head": {
            "seq": 58894,
            "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
            "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
            "timestamp": 1537581604,
            "fee": 485194,
            "version": 0,
            "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
            "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
        },
        "unspents": 38171,
        "unconfirmed": 1,
        "time_since_last_block": "4m46s"
    },
    "version": {
        "version": "0.24.1",
        "commit": "8798b5ee43c7ce43b9b75d57a1a6cd2c1295cd1e",
        "branch": "develop"
    },
    "open_connections": 8,
    "uptime": "6m30.629057248s",
    "csrf_enabled": true,
    "csp_enabled": true,
    "wallet_api_enabled": true,
    "gui_enabled": true,
    "unversioned_api_enabled": false,
    "json_rpc_enabled": false
}
`

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(jsonResponse)); err != nil {
		log.WithError(err).Panicln("unable to write response")
	}
}

func writeListenPortInFile(port uint16) (err error) {
	path := testrand.FilePathToSharePort()
	if !file.ExistFile(path) {
		var file, err = os.Create(path)
		if err != nil {
			log.WithError(err).Errorln("error creating file")
			return err
		}
		defer file.Close()
	}
	var file *os.File
	file, err = os.OpenFile(path, os.O_WRONLY, 0400)
	if err != nil {
		log.WithError(err).Errorln("error opening file")
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", port))
	if err != nil {
		log.WithError(err).Errorln("error writing file")
		return err
	}
	err = file.Sync()
	if err != nil {
		log.WithError(err).Errorln("error flushing file")
		return err
	}
	return err
}

func main() {
	var fakeNodePort = testrand.RandomPort()
	if err := writeListenPortInFile(fakeNodePort); err != nil {
		log.Fatal(err)
	}
	handler := http.HandlerFunc(httpHandler)
	log.WithError(http.ListenAndServe(fmt.Sprintf(":%d", fakeNodePort), handler)).Fatalln("server fail")
}
