package h_boot

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"net"
	"net/http"
	"time"
)

type RedirectHttps struct{}

func (*RedirectHttps) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "https"
	if r.Host == "" {
		r.URL.Host = "localhost"
	} else {
		if h, _, err := net.SplitHostPort(r.Host); err == nil {
			r.URL.Host = h
		} else {
			r.URL.Host = r.Host
		}
	}
	http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
}

//开启https后需要开启80端口做非https转发到https
func NewRedirectHttps() *http.Server {
	server := &http.Server{
		Addr:         a_boot.Config.App.Server.IP + ":80",
		WriteTimeout: time.Millisecond * 300,
		ReadTimeout:  time.Millisecond * 500,
		Handler:      &RedirectHttps{},
	}
	return server
}
