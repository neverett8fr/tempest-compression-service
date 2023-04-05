package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"tempest-compression-service/pkg/infra/compression"

	"github.com/gorilla/mux"
)

func newCompressionOperation(r *mux.Router) {
	r.HandleFunc("/compression/decompress", decompress).Methods(http.MethodPost)
	r.HandleFunc("/compression/compress", compress).Methods(http.MethodPost)
}

func decompress(w http.ResponseWriter, r *http.Request) {

	// post request from data-service
	// with body in bytes

	// Read the request body into a byte slice
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeFile(w, []byte("error reading body"))
		return
	}
	defer r.Body.Close()

	decompressedBody, err := CompressionProvider.Decompress(body)
	if err != nil {
		writeFile(w, []byte("error decompressing"))
		return
	}

	// Set the appropriate content type and content length headers
	w.Header().Set("Content-Type", compression.GetFileType(decompressedBody))
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(decompressedBody)))

	writeFile(w, decompressedBody)
}

func compress(w http.ResponseWriter, r *http.Request) {

	// post request from data-service
	// with body in bytes

	// Read the request body into a byte slice
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeFile(w, []byte("error reading body"))
		return
	}
	defer r.Body.Close()

	compressedBody, err := CompressionProvider.Compress(body)
	if err != nil {
		writeFile(w, []byte("error compressing"))
		return
	}

	// Set the appropriate content type and content length headers
	w.Header().Set("Content-Type", compression.GetFileType(compressedBody))
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(compressedBody)))

	writeFile(w, compressedBody)
}
