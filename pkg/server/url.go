package server

import (
	"errors"
	"net/http"

	"github.com/gustavocd/url-shortener/models"
	"github.com/spf13/viper"

	validation "github.com/go-ozzo/ozzo-validation"
	shortid "github.com/ventu-io/go-shortid"
)

var (
	// ErrInvalidCode occurs when the given code is not valid
	ErrInvalidCode = errors.New("La URL enviada es inválida")
	// ErrNotFound occurs when there is no related data to the query
	ErrNotFound = errors.New("No se encontraron resultados")
	// ErrCreateURL occurs when the create action fails
	ErrCreateURL = errors.New("No se pudo crear la URL")
	// ErrDecodeURL occurs when the given data can not be decoded
	ErrDecodeURL = errors.New("No es posible leer la información enviada")
	// ErrGenerateShortID occurs when we can't generate a unique short ID
	ErrGenerateShortID = errors.New("No es posible generar la URL corta, intentelo otra vez")
)

// HandleURLCreate handles URL creation
func (s *Server) HandleURLCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := &models.URL{}
		err := s.decode(w, r, url)
		if err != nil {
			errs := validation.Errors{}
			errs["message"] = ErrDecodeURL
			s.respondWithErr(w, r, errs, http.StatusBadRequest)
			return
		}

		err = url.Validate()
		if err != nil {
			s.respondWithErr(w, r, err, http.StatusBadRequest)
			return
		}

		sid, err := shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			errs := validation.Errors{}
			errs["message"] = ErrGenerateShortID
			s.respondWithErr(w, r, errs, http.StatusBadRequest)
			return
		}

		urlCode, err := sid.Generate()
		if err != nil {
			errs := validation.Errors{}
			errs["message"] = ErrGenerateShortID
			s.respondWithErr(w, r, errs, http.StatusBadRequest)
			return
		}

		url.URLCode = urlCode
		url.ShortURL = viper.GetString("BASE_URL") + urlCode

		err = s.db.Create(url)
		if err != nil {
			errs := validation.Errors{}
			errs["message"] = ErrCreateURL
			s.respondWithErr(w, r, errs, http.StatusInternalServerError)
			return
		}

		response := make(map[string]string, 2)
		response["message"] = "URL creada exitosamente"
		response["short_url"] = url.ShortURL

		s.respond(w, r, response, http.StatusCreated)
	}
}

// HandleURLRedirect handles URL redirect to the original URL
func (s *Server) HandleURLRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := s.parseCode(r)
		if err != nil {
			errs := validation.Errors{}
			errs["message"] = ErrInvalidCode
			s.respondWithErr(w, r, errs, http.StatusBadRequest)
			return
		}

		url := &models.URL{URLCode: code}
		err = s.db.Where("url_code = ?", code).First(url)
		if err != nil {
			errs := validation.Errors{}
			errs["message"] = ErrNotFound
			s.respondWithErr(w, r, errs, http.StatusNotFound)
			return
		}

		http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
	}
}
