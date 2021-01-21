package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const SearchURI = "https://www.data.qld.gov.au/api/3/action/package_search?q=organization:griffith-university"

type SearchResults struct {
	Help   string       `json:"help"`
	Success bool `json:"success"`
	Result SearchResult `json:"result"`
}

type SearchResult struct {
	Count   int      `json:"count"`
	Results []Result `json:"results`
}

type Result struct {
	Author        string     `json:"author"`
	AuthorEmail   string     `json:"author_email"`
	CreatorUserId string     `json:"creator_user_id"`
	Id            string     `json:"id"`
	Maintainer    string     `json:"maintainer"`
	Name          string     `json:"name"`
	Resources     []Resource `json:"resources"`
}

type Resource struct {
	Created      string `json:"created"`
	Description  string `json:"description"`
	Format       string `json:"format"`
	Id           string `json:"id"`
	LastModified string `json:"last_modified"`
	Mimetype     string `json:"mimetype"`
	Name         string `json:"name"`
	OriginalUrl  string `json:"original_url"`
	RevisionId   string `json:"revision_id"`
	Url          string `json:"url"`
}

func in(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func dieOnError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadResource(uri string) {
	colors := []string{"\x1b[37m", "\x1b[90m"}

	fmt.Printf("Loading: %s\n\n", uri)

	res, err := http.Get(uri)
	dieOnError(err)

	reader := csv.NewReader(res.Body)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	dieOnError(err)

	found_header := false
	for i, record := range records {
		if !found_header {
			found_header = in(record, "Description")
			if !found_header {
				continue
			}
		}
		fmt.Printf("%s%s\x1b[0m", colors[i%2], strings.Join(record, "\t"))
	}
}

func main() {
	res, err := http.Get(SearchURI)
	dieOnError(err)

	var data_raw SearchResults
	err = json.NewDecoder(res.Body).Decode(&data_raw)
	dieOnError(err)

	for _, result := range data_raw.Result.Results {
		if strings.Contains(result.Name, "contract") {
			for _, res := range result.Resources {
				loadResource(res.Url)
			}
		}
	}
}
