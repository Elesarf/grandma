package main

import (
	"errors"
	"os"
)

// ConfigData contains config data. Api keys, paths
type ConfigData struct {
	botAPIKey string
	picPath   string
	dbPath    string
}

func loadConfigFromEnv() (ConfigData, error) {

	botAPIKey, isSet := os.LookupEnv("GRANDMA_TELEGRAMM_BOT_API_KEY")
	if !isSet {
		return ConfigData{}, errors.New("No bot api key found. Please set GRANDMA_TELEGRAMM_BOT_API_KEY env")
	}

	picPath, isSet := os.LookupEnv("GRANDMA_PICTURE_PATH")
	if !isSet {
		return ConfigData{}, errors.New("No grandma picture path found. Please set GRANDMA_PICTURE_PATH env")
	}

	dbPath, isSet := os.LookupEnv("GRANDMA_DB_PATH")
	if !isSet {
		return ConfigData{}, errors.New("Error: not found picture directory. Please set GRANDMA_DB_PATH env")
	}

	return ConfigData{botAPIKey, picPath, dbPath}, nil
}
