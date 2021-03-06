package controllers

import (
	"net/http"

	schema "github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true) // ignore CSRF token key
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
