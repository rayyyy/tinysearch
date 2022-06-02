package tinysearch

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Engine struct {
	tokenizer     *Tokenizer
	indexer       *Indexer
	documentStore *DocumentStore
	indexDir      string
}

func NewSearchEngine(db *sql.DB) *Engine {
	tokenizer := NewTokenizer()
	indexer := NewIndexer(tokenizer)
	documentStore := NewDocumentStore(db)

	path, ok := os.LookupEnv("INDEX_DIR_PATH")
	if !ok {
		current, _ := os.Getwd()
		path = filepath.Join(current, "_index_data")
	}

	return &Engine{tokenizer, indexer, documentStore, path}
}

func (e *Engine) AddDocument(title string, reader io.Reader) error {
	id, err := e.documentStore.save(title)

	if err != nil {
		return err
	}
	e.indexer.update(id, reader)
	return nil
}

func (e *Engine) Flush() error {
	writer := NewIndexWriter(e.indexDir)
	return writer.Flush(e.indexer.index)
}

// 検索を実行する
func (e *Engine) Search(query string, k int) ([]*SearchResult, error) {

	// クエリをトークンに分割
	terms := e.tokenizer.TextToWordSequence(query)

	// 検索を実行
	docs := NewSearcher(e.indexDir).SearchTopK(terms, k)

	// タイトルを取得
	results := make([]*SearchResult, 0, k)
	for _, result := range docs.scoreDocs {
		title, err := e.documentStore.fetchTitle(result.docID)
		if err != nil {
			return nil, err
		}
		results = append(results, &SearchResult{
			result.docID, result.score, title,
		})
	}
	return results, nil
}

// 検索結果を格納する構造体
type SearchResult struct {
	DocID DocumentID
	Score float64
	Title string
}

// String print SearchTopK result info
func (r *SearchResult) String() string {
	return fmt.Sprintf("{DocID: %v, Score: %v, Title: %v}",
		r.DocID, r.Score, r.Title)
}
