package main

import (
	"fmt"
	"log"
	"os"

	hashindexes "github.com/phirmware/file-algos/hash-indexes"
)

func main() {
	applogger := log.New(os.Stdin, fmt.Sprintf("[%d]: Application: ", os.Getpid()), 2)
	db := hashindexes.NewDB()
	defer db.ShutDown()

	done := db.Write("/student/name", "phirmware")
	if done {
		value := db.Read("/student/name")
		applogger.Println(value)
	}
	// db.FlushDB()
}
