package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/frontend"
	"github.com/buexplain/go-fool"
)

func frontend(mux *fool.Mux) {
	mux.Get("/", c_frontend.Index)
	mux.Get("/index-widget", c_frontend.IndexWidget).AddLabel("json")
}
