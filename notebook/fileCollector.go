package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var oldTimeHours int64 = 6

// (?:[\D][a-zA-Z0-9]+[_]){3}([0-9]{10}){1}(\.)(ext|qne)[\W\D]

func removeOldFiles(lowTimeBound int64, db *sql.DB) {
	log.Println("Process remove old files")
	files := getOldestPicPathes(lowTimeBound, db)
	if len(files) > 1000 { // some magic
		log.Println("Remove " + strconv.FormatInt(int64(len(files)), 10) + " files")
		for _, file := range files {
			os.Remove(file)
		}
	}
}

func walidatePaths(db *sql.DB) {
	log.Println("Process validate paths")
	paths := getPicturesPath(db)
	for _, path := range paths {
		if !FileExists(path) {
			removeFileOnPath(path, db)
		}
	}
}

func checkFiles(path string, dbPath string) {

	var wg sync.WaitGroup
	log.Println("Process picture directory")
	log.Println("Directory is: " + path)

	re, err := regexp.Compile("(?:[\\D][a-zA-Z0-9]+[_]){3}([0-9]{10}){1}(\\.)(jpg|png|jpeg)")
	if err != nil {
		log.Panic(err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	oldTime := time.Now().Unix() + oldTimeHours*60*60
	removeOldFiles(oldTime, db)

	walidatePaths(db)

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if re.MatchString(file.Name()) {
			nameWithoutExt := strings.Split(file.Name(), ".")[0]
			fields := strings.Split(nameWithoutExt, "_")
			timestamp, _ := strconv.ParseUint(fields[3], 10, 64)
			wg.Add(1)
			go addToDB(PhotoForDB{fields[0], fields[1], fields[2], timestamp, path + "/" + file.Name()}, db, &wg)
		}
	}

	wg.Wait()
}
