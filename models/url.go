package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// URL represents the info related to a shorter URL
type URL struct {
	ID        int64     `json:"id" db:"id"`
	URLCode   string    `json:"url_code" db:"url_code"`
	LongURL   string    `json:"long_url" db:"long_url"`
	ShortURL  string    `json:"short_url" db:"short_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Validate validates URL data
func (u URL) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.LongURL,
			validation.Required.Error("El campo URL es obligatorio"),
			is.URL.Error("El formato de la URL no es válido"),
		),
	)
}
