package main

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
)

// SELECT NEWEST PIC: select path, timeStamp from pictures where userName="UserName" order by timeStamp DESC limit 1

func addToDB(photo PhotoForDB, db *sql.DB, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Print(photo.Name + " " + photo.MainLocation + " " + strconv.FormatUint(photo.TimeStamp, 10) + " " + photo.Path)

	log.Println("Insert into db")

	_, err := db.Exec(
		"insert into pictures (userName, mainLocation, addLocation, timestamp, path) values ($1, $2, $3, $4, $5)",
		photo.Name, photo.MainLocation, photo.AddLocation, photo.TimeStamp, photo.Path)
	return err
}

func getOldestPicPathes(lowTimeBound int64, db *sql.DB) []string {
	log.Println("Process get oldest pic paths")
	rows, err := db.Query("select path from pictures where timeStamp<$1", lowTimeBound)
	_check(err)
	defer rows.Close()

	pathes := []string{}
	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		_check(err)
		pathes = append(pathes, path)
	}

	return pathes
}

func getPicturesPath(db *sql.DB) (paths []string) {
	rows, err := db.Query("select path from pictures")
	if err != nil {
		log.Panic(err)
	}

	pathes := []string{}
	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		_check(err)
		pathes = append(pathes, path)
	}
	return pathes
}

func removeFileOnPath(path string, db *sql.DB) error {
	_, err := db.Exec(
		"delete from pictures where path=$1", path)
	return err
}

func removeOld(lowTimeBound uint64, db *sql.DB) error {
	_, err := db.Exec("delete from pictures where timeStamp<" + strconv.FormatUint(lowTimeBound, 10))
	return err
}
