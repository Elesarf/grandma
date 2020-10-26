package main

import (
	"database/sql"
	"log"
)

func getEventsListFromNameAndLocation(userName string, location InternalLocation) (paths []string) {
	SystemConfig, err := loadConfigFromEnv()
	_check(err)

	log.Println("Open db path: " + SystemConfig.dbPath)
	log.Println("User name: " + userName)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select path, timestamp from pictures where userName=$1 and mainLocation=$2 and addLocation=$3 order by timeStamp DESC limit 100",
		userName,
		location.mainLocation,
		location.optionalLocation)

	_check(err)
	defer rows.Close()

	var path string
	var time int64
	for rows.Next() {
		err := rows.Scan(&path, &time)
		_check(err)
		paths = append(paths, path)
	}
	return paths
}

func getLocations(userName string) (locations map[string]string) {
	SystemConfig, err := loadConfigFromEnv()
	_check(err)

	log.Println("Open db path: " + SystemConfig.dbPath)
	log.Println("User name: " + userName)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select mainLocation, addLocation from pictures where userName=$1",
		userName)
	_check(err)
	defer rows.Close()

	locations = make(map[string]string)
	for rows.Next() {
		var mainLoc string
		var addLoc string
		err := rows.Scan(&mainLoc, &addLoc)
		_check(err)

		log.Print(mainLoc + " : " + addLoc)

		_, ok := locations[mainLoc]
		if !ok {
			locations[mainLoc] = addLoc
		}
	}
	return locations
}

func getNewestPath(userName string, location InternalLocation) (path string, timestamp int64) {
	SystemConfig, err := loadConfigFromEnv() // TODO: rewrite get dbpath
	_check(err)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	// select path from pictures where userName=$1 and mainLocation=$2  and addLocation=$1 order by timeStamp DESC limit 1"
	rows, err := db.Query("select path, timestamp from pictures where userName=$1 and mainLocation=$2 and addLocation=$3 order by timeStamp DESC limit 1",
		userName,
		location.mainLocation,
		location.optionalLocation)

	_check(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&path, &timestamp)
		_check(err)
	}

	return path, timestamp
}

func addFriendToDb(userName string, friendName string) error {
	SystemConfig, err := loadConfigFromEnv() // TODO: rewrite get dbpath
	_check(err)

	log.Println("Add friend into base (" + SystemConfig.dbPath + ")")
	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Println("Error open base")
		return err
	}
	defer db.Close()

	_, errDbInsert := db.Exec("insert into friends (userName, friend) values ($1, $2)",
		userName, friendName)
	if errDbInsert == nil {
		log.Println("Correctly add friend (" + friendName + ") to " + userName)
	}

	return errDbInsert
}

func removeFriendFromDb(userName string, friendName string) error {
	SystemConfig, err := loadConfigFromEnv() // TODO: rewrite get dbpath
	_check(err)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, errDbInsert := db.Exec("delete from friends where userName=$1 and friend=$2",
		userName, friendName)
	return errDbInsert
}

func getFriendList(userName string) (friends []string) {
	SystemConfig, err := loadConfigFromEnv() // TODO: rewrite get dbpath
	_check(err)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select friend from friends where userName=$1",
		userName)

	_check(err)
	defer rows.Close()

	var friend string
	for rows.Next() {
		err := rows.Scan(&friend)
		_check(err)
		friends = append(friends, friend)
	}
	return friends
}

func getIInFriendsOn(userName string) (friends []string) {
	SystemConfig, err := loadConfigFromEnv() // TODO: rewrite get dbpath
	_check(err)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select userName from friends where friend=$1",
		userName)

	_check(err)
	defer rows.Close()

	var friend string
	for rows.Next() {
		err := rows.Scan(&friend)
		_check(err)
		friends = append(friends, friend)
	}
	return friends
}

func thatUserIsFriend(userName string, askName string) bool {
	SystemConfig, err := loadConfigFromEnv() // TODO: rewrite get dbpath
	_check(err)

	db, err := sql.Open("sqlite3", SystemConfig.dbPath)
	if err != nil {
		log.Panic(err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("select userName from friends where friend=$1",
		userName)
	_check(err)

	var ask string
	for rows.Next() {
		err := rows.Scan(&ask)
		log.Println("Check: " + userName + " on " + askName)
		if err == nil {
			log.Println("Check " + askName)

			if ask == askName {
				return true
			}
		}
	}
	return false
}
