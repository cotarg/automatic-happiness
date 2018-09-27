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

type Connector interface {
	Fetch() ([]byte, error)
	Translate([]byte) ([]Object, error)
}

type VAMConnector struct{}

func (c Connector) Fetch() ([]byte, error) {
	// V&A museum API Docs: https://www.vam.ac.uk/api/
	vickyMuseumURL := "https://www.vam.ac.uk/api/json/museumobject/"

	resp, err := http.Get(vickyMuseumURL)
	if err != nil {
		return nil, fmt.Errorf("argh this is bad, NO INFO b/c: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't give you your info b/c: %v", err)
	}

	return body, nil
}

func (c Connector) Translate([]byte) ([]Object, error)

func main() {
	http.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {
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
