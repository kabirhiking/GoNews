package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

// Document represents a news article
type Document struct {
	ID      int
	Title   string
	Date    string
	Content string
}

// LoadCSV expects a CSV with header including: id,title,date,content
func LoadCSV(path string) ([]Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Read header
	_, err = r.Read()
	if err != nil {
		return nil, err
	}

	var docs []Document
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// support flexible CSV columns: try to map by position
		// assume columns: id,title,date,content
		id, _ := strconv.Atoi(rec[0])
		content := ""
		if len(rec) > 3 {
			content = rec[3]
		}
		var date string
		if len(rec) > 2 {
			date = rec[2]
		}
		var title string
		if len(rec) > 1 {
			title = rec[1]
		}
		docs = append(docs, Document{
			ID:      id,
			Title:   title,
			Date:    date,
			Content: content,
		})
	}
	return docs, nil
}


