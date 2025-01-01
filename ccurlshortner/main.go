package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
)

type UrlShortResp struct {
	Key      string `json:"key"`
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

type UrlRequest struct {
	Url string `json:"url" validate:"required,url"`
}

var (
	u         UrlRequest
	urls             = []UrlShortResp{}
	inputFile string = "urls.json"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{surl}", redirectUrl)
	mux.HandleFunc("POST /createurl", shortUrl)
	mux.HandleFunc("DELETE /{durl}", deleteUrl)

	fmt.Println("Server Listening to :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func redirectUrl(w http.ResponseWriter, r *http.Request) {
	urlkeyfound := false
	surl := r.PathValue("surl")

	urls := unmarshallfile(inputFile)
	for _, url := range urls {
		if url.Key == surl {
			http.Redirect(w, r, url.LongUrl, 200)
			urlkeyfound = true
			break
		}
	}

	if !urlkeyfound {
		http.Error(w, "URL not found", http.StatusBadRequest)
	}
}

func shortUrl(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	newHash := sha256.New()
	newHash.Write([]byte(u.Url))
	key := newHash.Sum(nil)
	shortUrl := fmt.Sprintf("http://localhost:8000/%x", key[:12])
	var nkey string = fmt.Sprintf("%x", key[:12])

	resp := &UrlShortResp{
		Key:      nkey,
		LongUrl:  u.Url,
		ShortUrl: shortUrl,
	}

	writeToFile(resp)
	urlResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(urlResp)
}

func deleteUrl(w http.ResponseWriter, r *http.Request) {
	durl := r.PathValue("durl")

	urls := unmarshallfile(inputFile)
	for ind, url := range urls {
		if url.Key == durl {
			urls = append(urls[:ind], urls[ind+1:]...)
			updatedUrls, err := json.MarshalIndent(urls, "", " ")
			if err != nil {
				log.Fatal(err)
			}
			os.WriteFile(inputFile, updatedUrls, 0644)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeToFile(data *UrlShortResp) {
	urls := unmarshallfile(inputFile)
	urls = append(urls, *data)
	updatedUrls, err := json.MarshalIndent(urls, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(inputFile, updatedUrls, 0644)
}

func unmarshallfile(f string) []UrlShortResp {
	file, err := os.ReadFile(f)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if len(file) > 0 {
		err = json.Unmarshal(file, &urls)
		if err != nil {
			log.Fatal(err)
		}
	}
	return urls
}
