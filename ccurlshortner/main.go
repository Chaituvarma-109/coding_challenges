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

	_ "modernc.org/sqlite"

	"github.com/go-playground/validator/v10"
	"github.com/urlshortner/urldb"
)

type UrlRequest struct {
	Url string `json:"url" validate:"required,url"`
}

var (
	u UrlRequest
	//go:embed schema.sql
	ddl     string
	queries *urldb.Queries
	ctx     context.Context
)

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

	mux.HandleFunc("/{surl}", redirectUrl)
	mux.HandleFunc("POST /createurl", shortUrl)
	mux.HandleFunc("DELETE /{durl}", deleteUrl)

	fmt.Println("Server Listening to :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func redirectUrl(w http.ResponseWriter, r *http.Request) {
	surl := r.URL.Path[1:]

	lurl, err := queries.SelectLongUrl(ctx, surl)
	if err != nil {
		http.Error(w, "URL not found", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	http.Redirect(w, r, lurl, http.StatusSeeOther)
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

func deleteUrl(w http.ResponseWriter, r *http.Request) {
	durl := r.URL.Path[1:]

	err := queries.DeleteUrl(ctx, durl)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
