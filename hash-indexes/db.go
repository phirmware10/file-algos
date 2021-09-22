package hashindexes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DB_FILE = "./logs"
	INDEX_FILE = "./index.json"
)

type (
	DB struct {
		f        *os.File
		byteSize int
	}

	HashTable map[string]HashIndexValue

	HashIndexValue struct {
		Location int
		Max      int
	}

)

func init() {
	pid := os.Getpid()
	log.SetPrefix(fmt.Sprintf("[%d]: DB: ", pid))
}

func NewDB() DB {
	f, err := os.OpenFile(DB_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	indexf, err := os.OpenFile(INDEX_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer indexf.Close()

	info, err := f.Stat()

	if err != nil {
		log.Fatalf("FATAL: Error opening db file, %s", err)
	}

	hf := DB{
		f:        f,
		byteSize: int(info.Size()),
	}

	return hf
}

func (hf DB) ShutDown() {
	log.Println("Shutting down DB")
	defer hf.f.Close()
}

func getIndexes() HashTable {
	b, err := os.ReadFile(INDEX_FILE)
	indexes := make(HashTable)

	json.Unmarshal(b, &indexes)
	if err != nil {
		log.Fatalf("FATAL: Could not open index file, %s", err)
	}
	return indexes
}

func (hf DB) index(key string, end int) {
	indexes := getIndexes()
	hiv := HashIndexValue{
		Location: hf.byteSize,
		Max: hf.byteSize + end,
	}
	indexes[key] = hiv

	fbyte, err := json.Marshal(indexes)
	if err != nil {
		log.Fatalf("FATAL: could not write to file, %s", err)
	}

	os.WriteFile(INDEX_FILE, fbyte, 0666)
}

func (hf DB) Write(key, value string) bool {
	formattedStr := key + ":" + value + "\n"
	nb, err := hf.f.Write([]byte(formattedStr))
	if err != nil {
		log.Printf("ERROR: Could not write to file, %s", err)
		return false
	}
	hf.index(key, nb)
	hf.byteSize = hf.byteSize + nb

	return true
}

func (hf DB) Read(key string) string {
	db, err := os.Open(DB_FILE)
	if err != nil {
		log.Fatalf("FATAL: Unable to open file, err : %s", err)
	}

	indexes := getIndexes()
	hiv := indexes[key]

	size := hiv.Max - hiv.Location
	data := make([]byte, size)
	db.ReadAt(data, int64(hiv.Location))

	sdata := string(data)
	splitdata := strings.Split(sdata, ":")
	return splitdata[len(splitdata) - 1]
}

func (hf DB) FlushDB() {
	os.Remove(INDEX_FILE)
	log.Println("Remove Index")
	os.Remove(DB_FILE)
	log.Println("Remove DB logs")
	log.Println("FLUSHED ALL")
}
