package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gorilla/mux"
)

// Server represents the application server
type Server struct {
	db     *pop.Connection
	Router *mux.Router
}

// Response ...
type Response struct {
	Data     interface{} `json:"data,omitempty"`
	HTTPCode int         `json:"http_code"`
}

// ResponseErr ...
type ResponseErr struct {
	HTTPCode int   `json:"http_code"`
	Error    error `json:"errors"`
}

// ServeHTTP is a method that allows my Server to be a http.Handle type
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Data:     data,
		HTTPCode: status,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatalf("could not encode data %v", err)
	}
}

func (s *Server) respondWithErr(w http.ResponseWriter, r *http.Request, errors error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := ResponseErr{
		HTTPCode: status,
		Error:    errors,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatalf("could not encode data %v", err)
	}
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *Server) parseCode(r *http.Request) (string, error) {
	code, ok := mux.Vars(r)["code"]
	if !ok {
		return "", ErrInvalidCode
	}
	return code, nil
}

// NewServer creates a Server data type
func NewServer(db *pop.Connection, r *mux.Router) Server {
	return Server{db: db, Router: r}
}
