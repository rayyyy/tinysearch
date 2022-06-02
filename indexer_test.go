package tinysearch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	collection := []string{
		"Do you quarrel, sir?",
		"Quarrel sir! no, sir!",
		"No better.",
		"Well, sir",
	}

	indexer := NewIndexer(NewTokenizer())

	for i, doc := range collection {
		indexer.update(DocumentID(i), strings.NewReader(doc))
	}

	actual := indexer.index
	expected := &Index{
		map[string]PostingsList{
			"better":  NewPostingsList(NewPosting(2, 1)),
			"do":      NewPostingsList(NewPosting(0, 0)),
			"no":      NewPostingsList(NewPosting(1, 2), NewPosting(2, 0)),
			"quarrel": NewPostingsList(NewPosting(0, 2), NewPosting(1, 0)),
			"sir":     NewPostingsList(NewPosting(0, 3), NewPosting(1, 1, 3), NewPosting(3, 1)),
			"well":    NewPostingsList(NewPosting(3, 0)),
			"you":     NewPostingsList(NewPosting(0, 1)),
		},
		4,
	}

	assert.Equal(t, expected, actual)
}
