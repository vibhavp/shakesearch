package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"regexp"

	"pulley.com/shakesearch/parse"
)

type API struct {
	Works []parse.Work
}

type response struct {
	Error   string       `json:"error,omitempty"`
	Results []workResult `json:"results"`
}

type workResult struct {
	Title   string  `json:"title"`
	Matches []match `json:"matches"`
}

type match struct {
	Lines     []string `json:"text"`
	MatchLine int      `json:"matchLine"`
}

type searchResult struct {
	indices []int
	title   string
}

func extractLines(text string, index int) ([]string, int) {
	// TODO: Make the number of lines extracted configurable
	var backIndex, forwardIndex int

	lines := make([]string, 3)
	lines[1], forwardIndex, backIndex = extractCurLine(text, index)
	lines[0], _, _ = extractLine(text, backIndex, true)
	lines[2], _, _ = extractLine(text, forwardIndex, false)

	return lines, 1
}

func extractCurLine(text string, index int) (string, int, int) {
	var line string
	forwardIndex := index + 1
	for {
		if forwardIndex >= len(text) {
			break
		}
		if text[forwardIndex] == '\n' {
			break
		}

		forwardIndex++
	}

	line = text[index:forwardIndex]

	backIndex := index
	for {
		if backIndex < 0 {
			break
		}
		if text[backIndex] == '\n' {
			break
		}
		backIndex--
	}

	line = text[backIndex+1:index] + line
	return strings.Trim(line, " "), forwardIndex, backIndex
}

func extractLine(text string, index int, backwards bool) (string, int, bool) {
	if backwards {
		index--
	} else {
		index++
	}

	start := index
	var last byte

	for {
		if index >= len(text) || index < 0 {
			var line string
			if backwards {
				line = text[index+1 : start+1]
			} else {
				line = text[start+1 : index]
			}

			return strings.Trim(line, " "), index, true
		}

		if backwards {
			index--
		} else {
			index++
		}

		if text[index] == '\n' {
			if last == '\n' {
				continue
			}

			var line string
			if backwards {
				line = text[index+1 : start+1]
			} else {
				line = text[start+1 : index]
			}

			return strings.Trim(line, " "), index, false
		}

		last = text[index]
	}
}

func (a *API) search(ctx context.Context, text string) []workResult {
	results := make([]workResult, len(a.Works))
	wait := new(sync.WaitGroup)

	wait.Add(len(a.Works))

	for i := range a.Works {
		go func(i int) {
			select {
			case <-ctx.Done():
			default:
				work := a.Works[i]
				indices := work.Index.FindAllIndex(regexp.MustCompile("(?i)"+regexp.QuoteMeta(text)), -1)

				results[i] = workResult{
					Title: work.Title,
				}

				results[i].Matches = make([]match, len(indices))

				// this could be further parallelized, but this
				// is where I realised that I was supposed to
				// think about it from the *user's perspective*,
				// and not a purely backend performance one
				for i2, index := range indices {
					var lines []string
					var matchLine int
					lines, matchLine = extractLines(string(work.Index.Bytes()), index[0])
					results[i].Matches[i2].Lines = lines
					results[i].Matches[i2].MatchLine = matchLine
				}
			}
			wait.Done()
		}(i)
	}

	wait.Wait()
	return results
}

func (a *API) Search(w http.ResponseWriter, r *http.Request) {
	query, ok := r.URL.Query()["search"]
	if !ok || len(query[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response{Error: "missing search query in URL params"})
		return
	}

	results := a.search(r.Context(), query[0])

	json.NewEncoder(w).Encode(&response{
		Results: results,
	})
}
