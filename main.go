package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

var clicks atomic.Int64

func main() {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /time", timeHandler)
	mux.HandleFunc("POST /click", clickHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := page("Cloudless").Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	if err := serverTime(time.Now()).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func clickHandler(w http.ResponseWriter, r *http.Request) {
	n := clicks.Add(1)
	if err := clickCount(n).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
