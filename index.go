package tinysearch

import (
	"container/list"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

type Index struct {
	Dictionary     map[string]PostingList
	TotalDocsCount int
}

func NewIndex() *Index {
	dict := make(map[string]PostingList)
	return &Index{dict, 0}
}

type DocumentID int64

type Posting struct {
	DocID        DocumentID
	Positions    []int
	TermFreqency int
}

func NewPosting(docID DocumentID, positions ...int) *Posting {
	return &Posting{docID, positions, len(positions)}
}

type PostingList struct {
	*list.List
}

func NewPostingList(postings ...*Posting) PostingList {
	l := list.New()
	for _, posting := range postings {
		l.PushBack(posting)
	}
	return PostingList{l}
}

func (pl PostingList) add(posting *Posting) {
	pl.PushBack(posting)
}

func (pl PostingList) last() *Posting {
	e := pl.List.Back()
	if e == nil {
		return nil
	}
	return e.Value.(*Posting)
}

func (pl PostingList) Add(new *Posting) {
	last := pl.last()
	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return
	}
	last.Positions = append(last.Positions, new.Positions...)
	last.TermFreqency++
}

func (idx Index) String() string {
	var padding int
	keys := make([]string, 0, len(idx.Dictionary))
	for k := range idx.Dictionary {
		l := utf8.RuneCountInString(k)
		if l > padding {
			padding = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := "%-" + string(padding) + "s: %s"
	for i, k := range keys {
		if postingList, ok := idx.Dictionary[k]; ok {
			strs[i] = fmt.Sprintf(format, k, postingList.String())
		}
	}
	return fmt.Sprintf("total documents : %v\ndictionary:\n%v\n", idx.TotalDocsCount, strings.Join(strs, "\n"))
}

func (pl PostingList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e = e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}
	return strings.Join(str, "=>")
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v,%v,%v)", p.DocID, p.TermFreqency, p.Positions)
}
