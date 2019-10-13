package config

import "net/http"

//session 配置
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
