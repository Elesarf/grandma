package main

import (
	"log"
	"os"
)

func _check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
