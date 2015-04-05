package main

import "github.com/blevesearch/bleve"

// The generic data structure containing the information that we want to
// index with ES. The source of the searchitem (be it xml, json, etc)
// is irrelevant, as it is abstracted by the Crawlers themselves.
//
// Note that the structure matches the response expected from TapirGo,
// so we should keep it the same as long as we expect pairity.
//
// Also note: These are just strings for now.
type SearchItem struct {
	Score       string `json:"_score"`
	Content     string `json:"content"`
	Link        string `json:"link"`
	PublishedOn string `json:"published_on"`
	Summary     string `json:"summary"`
	Title       string `json:"title"`
}

// A searcher implements a simple interface for indexing and searching
// SearchItems.
type Searcher interface {
	// Index the SearchItem with whatever backend the Searcher is using.
	Index(...SearchItem) error

	// Search the backend for the given terms.
	Search(string) ([]SearchItem, error)
}

// TODO: Implement the elasticsearch searcher.
type ElasticSearcher struct {
}

// A super naive Bleve implementation of a Searcher
//
// TODO: Move BleveSearcher into a sub-package (since it won't be part
// of the final baku implementation)
type BleveSearcher struct {
	index bleve.Index
}

func NewBleveSearcher(p string) (*BleveSearcher, error) {
	i, err := bleve.New(p, bleve.NewIndexMapping())
	if err == bleve.ErrorIndexPathExists {
		i, err = bleve.Open(p)
	}
	if err != nil {
		return nil, err
	}

	return &BleveSearcher{
		index: i,
	}, nil
}

func (s *BleveSearcher) Index(sis ...SearchItem) error {
	for _, si := range sis {
		// Use the link as the index id.. in theory to reduce duplicate
		// entries.. hopefully.. in theory.. /handwave magic
		err := s.index.Index(si.Link, si)
		if err != nil {
			// Bailing on the error here is not the best thing to do,
			// i think i want to change how Searcher.Index handles errors..
			// unsure at the moment.
			return err
		}
	}
	return nil
}

func (s *BleveSearcher) Search(q string) ([]SearchItem, error) {
	request := bleve.NewSearchRequest(bleve.NewMatchQuery(q))
	request.Fields = []string{
		"content", "title", "link",
	}
	r, err := s.index.Search(request)
	if err != nil {
		return nil, err
	}

	sis := make([]SearchItem, r.Hits.Len())
	for i, hit := range r.Hits {
		sis[i] = SearchItem{}
		// This area is a bit ugly, needs to be cleaned up.
		if s, ok := hit.Fields["content"].(string); ok {
			sis[i].Content = s
			sis[i].Summary = s
		}
		if s, ok := hit.Fields["link"].(string); ok {
			sis[i].Link = s
		}
		if s, ok := hit.Fields["title"].(string); ok {
			sis[i].Title = s
		}
	}

	return sis, nil
}
