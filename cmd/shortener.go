package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gorilla/mux"
	"github.com/gustavocd/url-shortener/configs"
	"github.com/gustavocd/url-shortener/pkg/server"
	"github.com/spf13/viper"
)

func main() {
	conn, err := pop.Connect("development")
	if err != nil {
		log.Fatalf("could not connect to the database %v", err)
		return
	}

	configs.LoadConfig()

	r := mux.NewRouter()
	svr := server.NewServer(conn, r)
	svr.Router.HandleFunc("/{code}", svr.HandleURLRedirect()).Methods("GET")
	svr.Router.HandleFunc("/api/v1/url/shorten", svr.HandleURLCreate()).Methods("POST")

	log.Println(fmt.Sprintf("Server running on http://localhost%s 🐹", viper.GetString("PORT")))
	err = http.ListenAndServe(viper.GetString("PORT"), svr.Router)
	if err != nil {
		log.Fatalf("could not run the server %v", err)
		return
	}
}
