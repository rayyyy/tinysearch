package tinysearch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func setup() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(db:3306)/tinysearch")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("TRUNCATE TABLE documents")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll("_index_data"); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir("_index_data", 0777); err != nil {
		log.Fatal(err)
	}
	return db
}

func TestMain(m *testing.M) {
	testDB = setup()
	defer testDB.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateIndex(t *testing.T) {
	engine := NewSearchEngine(testDB)

	type testDoc struct {
		title string
		body  string
	}

	docs := []testDoc{
		{"test1", "Do you quarrel, sir?"},
		{"test2", "No better."},
		{"test3", "Quarrel sir! no, sir!"},
	}

	for _, doc := range docs {
		r := strings.NewReader(doc.body)
		if err := engine.AddDocument(doc.title, r); err != nil {
			t.Fatalf("failed to add document %s: %v", doc.title, err)
		}
	}

	if err := engine.Flush(); err != nil {
		t.Fatalf("failed to flush index: %v", err)
	}

	type testCase struct {
		file        string
		postingsStr string
	}

	testCases := []testCase{
		{
			"_index_data/better",
			`[{"DocID":2,"Positions":[1],"TermFreqency":1}]`,
		},
		{
			"_index_data/no",
			`[{"DocID":2,"Positions":[0],"TermFreqency":1},
			 {"DocID":3,"Positions":[2],"TermFreqency":1}]`,
		},
		{
			"_index_data/do",
			`[{"DocID":1,"Positions":[0],"TermFreqency":1}]`,
		},
		{
			"_index_data/quarrel",
			`[{"DocID":1,"Positions":[2],"TermFreqency":1},
			 {"DocID":3,"Positions":[0],"TermFreqency":1}]`,
		},
		{
			"_index_data/sir",
			`[{"DocID":1,"Positions":[3],"TermFreqency":1},
			 {"DocID":3,"Positions":[1,3],"TermFreqency":2}]`,
		},
		{
			"_index_data/you",
			`[{"DocID":1,"Positions":[1],"TermFreqency":1}]`,
		},
	}

	for _, testCase := range testCases {
		func() {
			file, err := os.Open(testCase.file)
			if err != nil {
				t.Fatalf("failed to open file %s: %v", testCase.file, err)
			}
			defer file.Close()

			b, err := ioutil.ReadAll(file)
			if err != nil {
				t.Fatalf("failed to read file %s: %v", testCase.file, err)
			}

			got := string(b)
			var buf bytes.Buffer
			_ = json.Compact(&buf, []byte(testCase.postingsStr))
			want := buf.String()
			assert.Equal(t, want, got, "file %s", testCase.file)
		}()
	}
}
