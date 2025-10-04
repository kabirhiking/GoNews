package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	path := flag.String("p", "data/news.csv", "path to news CSV file")
	query := flag.String("q", "", "search query")
	limit := flag.Int("n", 10, "max results to show")
	stem := flag.Bool("stem", false, "enable stemming (optional)")
	flag.Parse()

	start := time.Now()
	docs, err := LoadCSV(*path)
	if err != nil {
		log.Fatalf("failed to load dataset: %v", err)
	}
	fmt.Printf("Loaded %d docs from %s in %v", len(docs), *path, time.Since(start))

	// enable stemming option (analyze.go will honor this variable)
	EnableStemming = *stem

	idxStart := time.Now()
	idx := NewIndex()
	for _, d := range docs {
		idx.AddDocument(d)
	}
	fmt.Printf("Indexed %d docs in %v", idx.N, time.Since(idxStart))

	if *query == "" {
		fmt.Println("No query provided. Use -q \"your query\"")
		return
	}

	searchStart := time.Now()
	results := idx.Search(*query)
	fmt.Printf("Search completed in %v — %d results", time.Since(searchStart), len(results))

	// show top results
	count := 0
	for _, r := range results {
		if count >= *limit {
			break
		}
		d := idx.Docs[r.DocID]
		snippet := MakeSnippet(d.Content, r.MatchedTerms)
		fmt.Printf("[%s] %s (score: %.4f)%s", d.Date, d.Title, r.Score, snippet)
		count++
	}
}