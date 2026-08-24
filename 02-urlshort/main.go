package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/shroukelzoghby/gophercises_exercises/urlshort/urlshort"
)

func main() {

	flagYaml := flag.String("yaml", "yamlData.yaml", "Path to YAML file containing URL mappings")
	flagJSON := flag.String("json", "jsonData.json", "Path to JSON file containing URL mappings")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls, err := loadFromDB()
	if err != nil {
		fmt.Println("Error loading URL mappings from database:", err)
		return
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	//   - path: /urlshort
	//     url: https://github.com/gophercises/urlshort
	//   - path: /urlshort-final
	//     url: https://github.com/gophercises/urlshort/tree/solution
	//   `

	yaml, err := os.ReadFile(*flagYaml)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		fmt.Println("Error creating YAML handler:", err)
		return
	}

	// Build the JSONHandler using the yamlHandler as the
	// fallback
	jsonData, err := os.ReadFile(*flagJSON)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	jsonHandler, err := urlshort.JSONHandler(jsonData, yamlHandler)
	if err != nil {
		fmt.Println("Error creating JSON handler:", err)
		return
	}
	fmt.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", jsonHandler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func loadFromDB() (map[string]string, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/url?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT path, url FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pathsToUrls := make(map[string]string)
	for rows.Next() {
		var path string
		var url string
		if err := rows.Scan(&path, &url); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		pathsToUrls[path] = url
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pathsToUrls, nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
