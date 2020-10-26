package main

import (
	"log"
	"os"
)

// ConfigData contains system config
type ConfigData struct {
	pictureSavePath string
}

func loadConfigFromEnv() ConfigData {

	picturePath, isSet := os.LookupEnv("WATCHER_PICTURE_PATH")
	if !isSet {
		log.Panic("No picture path found. Please set WATCHER_PICTURE_PATH")
	}

	log.Println(picturePath)
	return ConfigData{picturePath}
}
