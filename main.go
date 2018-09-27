package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Object struct {
	ID                        string
	Title                     string
	Creator                   string
	Collection                string
	Completion                time.Time
	ColonialAppropriationDate time.Time
}

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, museumObjectQuery())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func museumObjectQuery() string {
	// V&A museum API Docs: https://www.vam.ac.uk/api/
	vickyMuseumURL := "https://www.vam.ac.uk/api/json/museumobject/"

	resp, err := http.Get(vickyMuseumURL)
	if err != nil {
		log.Fatal("OH MY GOSH THIS IS SO BROKEN YOU'RE PROBABLY SAD!!!")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("OH MY GOSH THIS IS SO BROKEN BUT IT'S DIFFERENT BROKEN! 💯")
	}

	return string(body)
}
