package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
	checks and logs errors
*/
func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

/*
	Reads the YAML/JSON files and returns the corresponding
	handlers
*/
func getHandler(fileName *string, mode *string, mapHandler http.HandlerFunc) http.HandlerFunc {
	*fileName = *fileName + "." + *mode
	fileContent, err := ioutil.ReadFile(*fileName)
	check(err)
	if *mode == "yaml" {
		tempHandler, err := YAMLHandler([]byte(fileContent), mapHandler)
		check(err)
		return tempHandler
	} else {
		tempHandler, err := JSONHandler([]byte(fileContent), mapHandler)
		check(err)
		return tempHandler
	}
}

func main() {
	/*
		Flags:
		mode :- select the mode (yaml,json, db)
		filename :- if mode is yaml/json, then name of the yaml/json file (excluding the extension)
		dbname :- if mode is db, name of the database
	*/
	mode := flag.String("mode", "yaml", "Give mode as input. Options :- yaml, json, db")
	fileName := flag.String("filename", "default", "Name of the YAML/JSON file. Specify only if mode is yaml/json")
	dbName := flag.String("dbname", "url_maps", "Name of the database. Specify only if mode is db")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	var netHandler http.HandlerFunc

	if *mode == "yaml" || *mode == "json" {
		if *dbName != "url_maps" {
			log.Fatalln("Incorrect combination of flags")
			os.Exit(1)
		}
		netHandler = getHandler(fileName, mode, mapHandler)

	} else if *mode == "db" {
		tempHandler, err := DBHandler(*dbName, mapHandler)
		check(err)
		netHandler = tempHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", netHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
