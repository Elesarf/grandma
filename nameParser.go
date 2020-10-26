package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// PictureValues keep picture
type PictureValues struct {
	UTCTime            time.Time
	mainLocation       string
	additionalLocation string
	userName           string
}

// CheckPictureName checked picture name
func CheckPictureName(name string) bool {
	return true
}

// ConstructName ...
func ConstructName(name string) (PictureValues, error) {

	entity := strings.Split(name, "_")

	if len(entity) != 4 {
		return PictureValues{}, errors.New("Wrong name string")
	}

	log.Println("Parse values:")
	log.Println("\tname: " + entity[0])
	log.Println("\tmain location: " + entity[1])
	log.Println("\taddd location: " + entity[2])
	log.Println("\ttime : " + entity[3])

	userName := entity[0]
	mainLocation := entity[1]
	additionalLocation := entity[2]
	timeValue, err := strconv.ParseInt(entity[3], 10, 64)
	if err != nil {
		return PictureValues{}, errors.New("Wrong time")
	}

	timeParse := time.Unix(timeValue, 0)

	return PictureValues{timeParse, mainLocation, additionalLocation, userName}, nil
}
