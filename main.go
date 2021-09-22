package main

import (
	"fmt"

	hashindexes "github.com/phirmware/file-algos/hash-indexes"
)


func main() {
	db := hashindexes.NewDB()
	defer db.ShutDown()

	done := db.Write("/second/key", "second value")
	if done {
		value := db.Read("/first/key")
		fmt.Println(value)
	}
}
