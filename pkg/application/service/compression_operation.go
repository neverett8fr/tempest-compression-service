package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func newCompressionOperation(r *mux.Router) {
	r.HandleFunc("/compression/decompress", decompress).Methods(http.MethodPost)
	r.HandleFunc("/compression/compress", compress).Methods(http.MethodPost)
}

func decompress(w http.ResponseWriter, r *http.Request) {

	// post request from data-service
	// with body in bytes
	var body []byte
	_ = json.NewDecoder(r.Body).Decode(&body)

	uncompressedBody, err := CompressionProvider.Decompress(body)
	if err != nil {
		writeFile(w, []byte("error decompressing"))
		return
	}

	writeFile(w, uncompressedBody)
}

func compress(w http.ResponseWriter, r *http.Request) {

	// post request from data-service
	// with body in bytes
	var body []byte
	_ = json.NewDecoder(r.Body).Decode(&body)

	uncompressedBody, err := CompressionProvider.Compress(body)
	if err != nil {
		writeFile(w, []byte("error decompressing"))
		return
	}

	writeFile(w, uncompressedBody)
}
