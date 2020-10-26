package main

import (
	"log"
	"os"
)

var filePath string
var dbPath string

func main() {

	filePath, isSet := os.LookupEnv("GRANDMA_PICTURE_PATH")
	if !isSet {
		log.Panic("Error: not found picture directory. Please set GRANDMA_PICTURE_PATH env.")
	}

	dbPath, isSet := os.LookupEnv("GRANDMA_DB_PATH")
	if !isSet {
		log.Panic("Error: not found base. Please set GRANDMA_DB_PATH env.")
	}

	log.Println("Start process dir")
	checkFiles(filePath, dbPath)
}
