package config

import "net/http"

type Session struct {
	Store   string
	Key     string
	Name    string
	Options struct {
		Path     string
		Domain   string
		MaxAge   int
		Secure   bool
		HttpOnly bool
		SameSite http.SameSite
	}
}
