package main

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"pulley.com/shakesearch/api"
	"pulley.com/shakesearch/parse"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	works, err := parse.GetWorks(searcher.SuffixArray)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	a := &api.API{Works: works}

	// http.HandleFunc("/search", handleSearch(searcher))
	http.HandleFunc("/search", a.Search)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	SuffixArray *suffixarray.Index
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.SuffixArray = suffixarray.New(dat)

	return nil
}
