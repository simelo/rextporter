package main

import (
	"github.com/denisacostaq/rextporter/config"
	"log"
	"encoding/json"
	"os"
)


func main() {
	conf := config.Config()
	if b, err := json.MarshalIndent(conf, "", " "); err != nil {
		log.Println("Error marshalling:", err)
	} else {
		os.Stdout.Write(b)
	}
}
