package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		url, ok := pathsToUrls[path]
		if !ok {
			fallback.ServeHTTP(w, req)
			return
		}
		http.Redirect(w, req, url, http.StatusFound)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yamlRecords, err := yamlParse(yml)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]string)

	for _, record := range yamlRecords {
		urlMap[record.Path] = record.Url
	}
	return MapHandler(urlMap, fallback), nil
}

/*
	JSONHandler receives the json file, parses it and creates a map
	of paths to urls. It then uses this map to call the MapHandler
	which returns a handler of type http.HandlerFunc. On encountering
	any error it returns nil and an error to be handled in main
*/
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	jsonRecords, err := jsonParse(json)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]string)

	for _, record := range jsonRecords {
		urlMap[record.Path] = record.Url
	}
	return MapHandler(urlMap, fallback), nil
}

/*
	DBHandler implements the logic to connect to the database,
	retrieve the stored mapping data and creates a map between
	paths and the corresponding urls. It then calls the MapHandler
	using this map and the fallback to return a handler of type
	http.HandlerFunc
*/

func DBHandler(dbName string, fallback http.Handler) (http.HandlerFunc, error) {
	dbName = "tempuser:crCf1cezlgSp5N4l@tcp(127.0.0.1:3306)/" + dbName
	db, err := sql.Open("mysql", dbName)

	if err != nil {
		return nil, err
	}

	sql := "Select * from records"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]string)
	var url string
	var path string

	for rows.Next() {
		err := rows.Scan(&path, &url)
		if err != nil {
			return nil, err
		}
		urlMap[path] = url
	}
	err = rows.Err()
	rows.Close()
	db.Close()
	if err != nil {
		return nil, err
	}

	return MapHandler(urlMap, fallback), nil
}

/*
	struct to read the parsed yaml files
*/
type StructYaml struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

/*
	struct to read the parsed json files
*/
type StructJson struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

/*
	Parses the yaml file and returns an array
	of structs of type StructYaml
*/
func yamlParse(yml []byte) ([]StructYaml, error) {
	var yamlRecords []StructYaml
	err := yaml.Unmarshal(yml, &yamlRecords)
	//fmt.Println(temp)
	return yamlRecords, err
}

/*
	Parses the json file and returns an array
	of structs of type StructJson
*/
func jsonParse(jsn []byte) ([]StructJson, error) {
	var jsonRecords []StructJson
	err := json.Unmarshal(jsn, &jsonRecords)
	return jsonRecords, err
}
