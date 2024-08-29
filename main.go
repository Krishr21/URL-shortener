package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var UrlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher:", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data:", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Encodetostring:", hash)
	fmt.Println("final string:", hash[:8])
	return hash[:8]
}

func CreateURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	id := shortURL
	UrlDB[id] = URL{
		ID:           id,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL

}

func getURL(id string) (URL, error) {
	url, ok := UrlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil

}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL_ := CreateURL(data.URL)
	//fmt.Fprintf(w,shortURL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)

}

func main() {
	//fmt.Println("Starting URL shortner...")
	//OriginalURL := "https://github.com/Krishr21"
	//generateShortURL(OriginalURL)

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)

	fmt.Println("starting seerver on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("on starting server", err)
	}
}
