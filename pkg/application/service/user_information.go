package service

import (
	"fmt"
	"net/http"
	application "tempest-user-service/pkg/application/entities"

	"github.com/gorilla/mux"
)

func NewUserInformation(r *mux.Router) {
	r.HandleFunc("/test/{text}", testHandler).Methods("GET")
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	text := params["text"]

	body := application.NewResponse(fmt.Sprintf("test: %v", text))

	writeReponse(w, r, body)
}
