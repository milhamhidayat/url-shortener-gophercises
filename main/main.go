package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"url-shortener-gophercises"

	"github.com/boltdb/bolt"
)

func main() {

	var yamlFile, jsonFile, boldDbFile string
	flag.StringVar(&yamlFile, "yaml", "", "YAML file for router")
	flag.StringVar(&jsonFile, "json", "", "JSON file for router")
	flag.StringVar(&boldDbFile, "boltdb", "", "Bolt DB file for router")

	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var urlData []byte
	var err error
	var urlHandler http.HandlerFunc

	if yamlFile != "" {
		urlData, err = ioutil.ReadFile(yamlFile)

		if err != nil {
			panic(err)
		}

		urlHandler, err = urlshort.YAMLHandler([]byte(urlData), mapHandler)

		if err != nil {
			panic(err)
		}
	} else if jsonFile != "" {
		urlData, err = ioutil.ReadFile(jsonFile)

		if err != nil {
			panic(err)
		}

		urlHandler, err = urlshort.JSONHandler([]byte(urlData), mapHandler)

		if err != nil {
			panic(err)
		}
	} else if boldDbFile != "" {
		insertRoute()

		db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			panic(err)
		}
		defer db.Close()

		urlHandler, err = urlshort.BoltDbHandler(db, mapHandler)

		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", urlHandler)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	// 	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mapHandler)
	// http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func insertRoute() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("urlBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte("/urlshort-godoc"), []byte("https://godoc.org/github.com/gophercises/urlshort"))
		err = b.Put([]byte("/yaml-godoc"), []byte("https://godoc.org/gopkg.in/yaml.v2"))

		return nil
	})
}
