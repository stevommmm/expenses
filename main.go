package main

import (
	"net/http"
	"encoding/json"
	"encoding/csv"
	"fmt"
	"strings"
	"log"
)

const SearchURI = "https://www.data.qld.gov.au/api/3/action/package_search?q=organization:griffith-university"

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

	fmt.Printf("%#v\n\n", uri)

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
		fmt.Printf("%s%s\x1b[0m", colors[i % 2], strings.Join(record, "\t"))
	}
}

// There are neater ways to handle this, rewrite for nested structs
func parseResources(raw map[string]interface{}) map[string][]string {
	results := map[string][]string{}

	if data_result, ok := raw["result"].(map[string]interface{}); ok {
		if data_results, ok := data_result["results"]; ok {
			for _, rinter := range data_results.([]interface{}) {
				result := rinter.(map[string]interface{})
				if resources, ok := result["resources"]; ok {
					for _, rres := range resources.([]interface{}) {
						resource := rres.(map[string]interface{})
						if url, ok := resource["url"]; ok {
							results[result["name"].(string)] = append(results[result["name"].(string)], url.(string))
						}
					}
				}
			}
		}
	}
	return results
}

func main() {
	res, err := http.Get(SearchURI)
	dieOnError(err)

	var data_raw map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data_raw)
	dieOnError(err)

	resources := parseResources(data_raw)

	for name, resources := range resources {
		if strings.Contains(name, "contract") {
			for _, uri := range resources {
				loadResource(uri)
			}
		}
	}
}
