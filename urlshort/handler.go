package urlshort

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type urlRedirect struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url"  json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirect, ok := pathsToUrls[path]

		if ok {
			http.Redirect(w, r, redirect, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	data := parseYaml([]byte(yml))
	dataMap := buildMap(data)
	mh := MapHandler(dataMap, fallback)

	return mh, nil
}

func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	data := parseJson(yml)
	dataMap := buildMap(data)
	mh := MapHandler(dataMap, fallback)

	return mh, nil
}

func SQLHandler(db *sql.DB, fallback http.Handler) (http.HandlerFunc, error) {
	rows, err := db.Query("SELECT * FROM entries")
	if err != nil {
		log.Fatalf("Error reading database: %v\n", err)
		os.Exit(1)
	}

	defer rows.Close()

	dataMap := make(map[string]string)
	for rows.Next() {
		var id int
		var path string
		var url string
		if err := rows.Scan(&id, &path, &url); err != nil {
			log.Fatalf("Error parsing database rows: %v\n", err)
		}
		dataMap[path] = url
	}

	mh := MapHandler(dataMap, fallback)

	return mh, nil
}

func parseYaml(yml []byte) []urlRedirect {
	var redirects []urlRedirect

	err := yaml.Unmarshal(yml, &redirects)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		os.Exit(1)
	}

	return redirects
}

func parseJson(yml []byte) []urlRedirect {
	var redirects []urlRedirect

	err := json.Unmarshal(yml, &redirects)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		os.Exit(1)
	}
	return redirects
}

func buildMap(ur []urlRedirect) map[string]string {
	m := make(map[string]string)

	for _, r := range ur {
		m[r.Path] = r.URL
	}
	return m
}
