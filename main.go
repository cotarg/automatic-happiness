package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MuseumObject struct {
	ID                        string     `json:"id"`
	Title                     string     `json:"title"`
	Creator                   string     `json:"creator"`
	Collection                string     `json:"collection"`
	Completion                *time.Time `json:"completion,omitempty"`
	ColonialAppropriationDate string     `json:"colonialAppropriationDate"`
}

type Connector interface {
	Fetch() ([]byte, error)
	Translate([]byte) ([]MuseumObject, error)
}

type VAMConnector struct{}

type VAMObjAPIResp struct {
	Records []struct {
		Fields struct {
			PrimaryImageID    string `json:"primary_image_id"`
			MuseumNumber      string `json:"museum_number"`
			Artist            string `json:"artist"`
			CollectionCode    string `json:"collection_code"`
			Location          string `json:"location"`
			DateText          string `json:"date_text"`
			MuseumNumberToken string `json:"museum_number_token"`
			Object            string `json:"object"`
			Longitude         string `json:"longitude"`
			ObjectNumber      string `json:"object_number"`
			Slug              string `json:"slug"`
			Latitude          string `json:"latitude"`
			Title             string `json:"title"`
			Place             string `json:"place"`
		} `json:"fields"`
		Pk    int    `json:"pk"`
		Model string `json:"model"`
	} `json:"records"`
	Meta struct {
		ResultCount int `json:"result_count"`
	} `json:"meta"`
}

func (c VAMConnector) Translate(body []byte) ([]MuseumObject, error) {
	objFromMuseum := &VAMObjAPIResp{}
	if err := json.Unmarshal(body, objFromMuseum); err != nil {
		return nil, fmt.Errorf("i can't do that Dave: %v", err)
	}

	museumObjects := make([]MuseumObject, 0)

	for _, obj := range objFromMuseum.Records {
		museumObjects = append(museumObjects, MuseumObject{
			ID:                        obj.Fields.MuseumNumber,
			Title:                     obj.Fields.Title,
			Creator:                   obj.Fields.Artist,
			Collection:                obj.Fields.CollectionCode,
			ColonialAppropriationDate: obj.Fields.DateText,
		})
	}
	return museumObjects, nil
}

func (c VAMConnector) Fetch() ([]byte, error) {
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

func main() {
	http.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {
		connector := VAMConnector{}
		body, err := connector.Fetch()
		if err != nil {
			http.Error(w, fmt.Sprintf("THIS FUNCTION FAILED! YOU HAVE NO %v", err), http.StatusTeapot)
		}

		resp, err := connector.Translate(body)
		if err != nil {
			http.Error(w, fmt.Sprintf("We don't have enough cool error codes so you get %v", err), http.StatusNotAcceptable)
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, fmt.Sprintf("Oh, you think you're getting JSON, huh? This is a YAML house! %v", err), http.StatusNoContent)
		}

		// I stole this next line from StackOverflow and I am not ashamed
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(jsonResp)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
