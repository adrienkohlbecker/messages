package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrienkohlbecker/messages/formatter"
	"github.com/adrienkohlbecker/messages/model"
)

var msgs model.Messages

func main() {

	files := []string{
		"/Users/ak/Desktop/messages-store/signal.json",
	}

	for _, file := range files {
		m, err := model.Load(file)
		if err != nil {
			log.Fatal(err)
		}
		msgs = append(msgs, m...)
	}

	http.HandleFunc("/", msgHandler)
	http.HandleFunc("/serve", serveHandler)

	log.Printf("Parsed %d messages", len(msgs))
	log.Printf("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func msgHandler(w http.ResponseWriter, r *http.Request) {

	err := formatter.Format(msgs, w)
	if err != nil {
		panic(err)
	}

}

func serveHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")

	if !strings.HasPrefix(path, "/Users/ak/Google Drive/Applications/Messages/media") {
		panic("unauthorized")
	}

	if strings.HasSuffix(path, ".jpg") {
		w.Header().Add("Content-Type", "image/jpeg")
	} else if strings.HasSuffix(path, ".mp4") {
		w.Header().Add("Content-Type", "video/mp4")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(500)
	}

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(500)
	}

}
