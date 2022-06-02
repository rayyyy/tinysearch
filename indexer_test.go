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

	indexer := NewIndex(NewTokenizer())

	for i, doc := range collection {
		indexer.Update(DocumentID(i), strings.NewReader(doc))
	}

	actual := indexer.index
	expected := &Index{
		map[string]PostingList{
			"better":  NewPostingList(NewPosting(2, 1)),
			"do":      NewPostingList(NewPosting(0, 0)),
			"no":      NewPostingList(NewPosting(1, 2), NewPosting(2, 0)),
			"quarrel": NewPostingList(NewPosting(0, 2), NewPosting(1, 0)),
			"sir":     NewPostingList(NewPosting(0, 3), NewPosting(1, 1, 3), NewPosting(3, 1)),
			"well":    NewPostingList(NewPosting(3, 0)),
			"you":     NewPostingList(NewPosting(0, 1)),
		},
		4,
	}

	assert.Equal(t, expected, actual)
}
