package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"tempest-compression-service/pkg/config"
	"tempest-compression-service/pkg/infra/compression"

	"github.com/gorilla/mux"
)

var (
	CompressionProvider compression.CompressionProvider
)

func NewRoutes(r *mux.Router, conf config.Config) {
	// initialise provider
	cp, err := compression.InitialiseCompressionProvider(
		context.Background(),
	)
	if err != nil {
		log.Printf("error initialising compression provider, err %v", err)
	}

	CompressionProvider = cp

	// initialise routes
	newCompressionInformation(r)
	newCompressionOperation(r)
}

func writeReponse(w http.ResponseWriter, r *http.Request, body interface{}) {

	reponseBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("error converting reponse to bytes, err %v", err)
	}
	w.Header().Add("Content-Type", "application/json")

	_, err = w.Write(reponseBody)
	if err != nil {
		log.Printf("error writing response, err %v", err)
	}
}
