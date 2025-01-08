package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/urlshortner/docs"
	"github.com/urlshortner/urldb"
	_ "modernc.org/sqlite"
)

// UrlRequest
//
//	@description	Request for creating shorturl
type UrlRequest struct {
	// url like "http://www.example.com", "https://www.google.com"
	Url string `json:"url" validate:"required,url"`
}

var (
	u UrlRequest
	//go:embed schema.sql
	ddl     string
	queries *urldb.Queries
	ctx     context.Context
)

// @title			Urlshortner API
// @version		1.0
// @description	This is a Url Shortner api.
// @termsOfService	http://swagger.io/terms/
func main() {
	ctx = context.Background()

	db, err := sql.Open("sqlite", "file:test.db")
	if err != nil {
		log.Fatal("cannot create db", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot ping database:", err)
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
		return
	}

	queries = urldb.New(db)

	mux := http.NewServeMux()

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/swagger/doc.json")))
	mux.HandleFunc("GET /{surl}", redirectUrl)
	mux.HandleFunc("POST /createurl", shortUrl)
	mux.HandleFunc("DELETE /{durl}", deleteUrl)

	fmt.Println("Server Listening to :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}

// redirectUrl
//
//	@Summary		Redirect Url
//	@Description	Redirect given short Url to original or long url
//	@Tags			Urls
//	@Accept			json
//	@Produce		json
//	@Param			surl	path		string		true	"url to redirect"
//	@Success		302		{string}	Location	"Redirects to the long url"
//	@Failure		404		{string}	http.StatusNotFound
//	@Router			/{surl} [get]
func redirectUrl(w http.ResponseWriter, r *http.Request) {
	surl := r.URL.Path[1:]

	lurl, err := queries.SelectLongUrl(ctx, surl)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
	}

	// Check if the request is coming from Swagger (or expects JSON)
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFound)
		return
	}

	http.Redirect(w, r, lurl, http.StatusFound)
}

// shortUrl
//
//	@Summary		Short Url
//	@Description	Create Short Url from Long or Original Url
//	@Tags			Urls
//	@Accept			json
//	@Produce		json
//	@Param			Url	body		UrlRequest	true	"this long or original url"
//	@Success		200	{object}	urldb.Url
//	@Failure		400	{string}	http.StatusBadRequest
//	@Failure		500	{string}	http.StatusInternalServerError
//	@Router			/createurl [post]
func shortUrl(w http.ResponseWriter, r *http.Request) {
	// Create a JSON decoder and disallow unknown fields
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the request into the UrlRequest struct
	if err := decoder.Decode(&u); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Check if the 'url' field is empty
	if u.Url == "" {
		http.Error(w, "The 'url' field is required and cannot be empty", http.StatusBadRequest)
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

	resp, err := queries.CreateUrl(ctx, urldb.CreateUrlParams{
		Key:      nkey,
		Longurl:  string(u.Url),
		Shorturl: shortUrl,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	urlResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(urlResp)
}

// deleteUrl
//
//	@Summary		Delete Url
//	@Description	Delete Url given its shorturl
//	@Tags			Urls
//	@Accept			json
//	@Produce		json
//	@Param			durl	path		string	true	"url to delete"
//	@Success		200		{string}	http.StatusOK
//	@Failure		404		{string}	http.StatusNotFound
//	@Router			/{durl} [delete]
func deleteUrl(w http.ResponseWriter, r *http.Request) {
	durl := r.URL.Path[1:]

	err := queries.DeleteUrl(ctx, durl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
