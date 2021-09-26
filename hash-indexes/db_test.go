package hashindexes

import (
	"testing"
)

func setupTestDB() DB {
	DB_FILE = "./testdata/logs"
	INDEX_FILE = "./testdata/index.json"

	return NewDB()
}

func TestNewDB(t *testing.T) {
	db := setupTestDB()
	fs, err := db.f.Stat()
	if err != nil {
		t.Errorf("ERROR: %s", err)
		return
	}

	got := fs.Name()
	expected := "logs"
	if expected != got {
		t.Errorf("Expected %s got %s", "logs", got)
	}

	expectedSize := 0
	gotSize := fs.Size()
	if expectedSize != int(gotSize) {
		t.Errorf("Expected %d got %d", expectedSize, gotSize)
	}
}

func TestGetIndexes(t *testing.T) {
	testtable := []struct {
		DB_FILE          string
		INDEX_FILE       string
		expectedLocation int
		expectedMax      int
	}{
		{
			DB_FILE:          "./testdata/log",
			INDEX_FILE:       "./testdata/index.json",
			expectedLocation: 45,
			expectedMax:      67,
		},
		{
			DB_FILE:          "./testdata/log",
			INDEX_FILE:       "./testdata/empty_index.json",
			expectedLocation: 0,
			expectedMax:      0,
		},
	}

	for _, table := range testtable {
		DB_FILE = table.DB_FILE
		INDEX_FILE = table.INDEX_FILE

		ht := getIndexes()
		gotLocation := ht["/unix/data"].Location
		gotMax := ht["/unix/data"].Max

		if table.expectedLocation != gotLocation {
			t.Errorf("Expected %d got %d", table.expectedLocation, gotLocation)
		}

		if table.expectedMax != gotMax {
			t.Errorf("Expected %d got %d", table.expectedMax, gotMax)
		}
	}
}
