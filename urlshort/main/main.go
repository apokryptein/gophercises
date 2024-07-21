package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/gophercises/urlshort"
	"urlshort"
)

type Flags struct {
	YamlFile string
	JsonFile string
	SqlDb    bool // if desired, can update this to take a connection string as opposed to hardcoding one
}

const (
	ServerPort = "8080"
)

func main() {
	flags := getFlags()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := createMapHandler(pathsToUrls)

	if flags.YamlFile != "" {
		yamlHandler := createYamlHandler(flags.YamlFile, mapHandler)
		serveHandler(ServerPort, yamlHandler)
	} else if flags.JsonFile != "" {
		jsonHandler := createJsonHandler(flags.JsonFile, mapHandler)
		serveHandler(ServerPort, jsonHandler)
	} else if flags.SqlDb {
		sqlHandler := createSqlHandler(mapHandler)
		serveHandler(ServerPort, sqlHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func getFlags() Flags {
	yamlFile := flag.String("y", "", "YAML file containing redirects")
	jsonFile := flag.String("j", "", "JSON file containing redirects")
	sqlDb := flag.Bool("db", false, "Use local PostgreSQL database")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Please select a file for input: JSON or YAML")
		flag.Usage()
		os.Exit(0)
	}

	flags := Flags{
		YamlFile: *yamlFile,
		JsonFile: *jsonFile,
		SqlDb:    *sqlDb,
	}

	return flags
}

func loadData(f string) ([]byte, error) {
	d, err := os.ReadFile(f)
	return d, err
}

func createYamlHandler(f string, fallback http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Loading YAML: ", f)
	data, err := loadData(f)
	if err != nil {
		log.Fatalf("Error loading file: %v\n", err)
		os.Exit(1)
	}

	yamlHandler, err := urlshort.YAMLHandler(data, fallback)
	if err != nil {
		log.Fatalf("Error creating handler: %v\n", err)
	}

	return yamlHandler
}

func createJsonHandler(f string, fallback http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Loading JSON: ", f)
	data, err := loadData(f)
	if err != nil {
		log.Fatalf("Error loading file: %v\n", err)
		os.Exit(1)
	}

	jsonHandler, err := urlshort.JSONHandler(data, fallback)
	if err != nil {
		log.Fatalf("Error loading JSON: %v\n", err)
	}

	return jsonHandler
}

func createSqlHandler(fallback http.HandlerFunc) http.HandlerFunc {
	connStr := "user=testuser dbname=redirects password=testpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error loading database: %v\n", err)
	}
	defer db.Close()

	sqlHandler, err := urlshort.SQLHandler(db, fallback)
	if err != nil {
		log.Fatalf("Error creating SQL Handler: %v\n", err)
	}

	return sqlHandler
}

func createMapHandler(u map[string]string) http.HandlerFunc {
	mux := defaultMux()
	mapHandler := urlshort.MapHandler(u, mux)
	return mapHandler
}

func serveHandler(p string, h http.Handler) {
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":"+p, h)
}
