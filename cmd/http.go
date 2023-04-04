package cmd

import (
	"fmt"
	"log"
	"net/http"
	"tempest-user-service/pkg/config"
	"time"

	"github.com/gorilla/mux"
)

func StartServer(conf *config.Service, router *mux.Router) error {
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server started on port: %v", conf.Port)

	log.Fatal(srv.ListenAndServe())
	return nil
}
