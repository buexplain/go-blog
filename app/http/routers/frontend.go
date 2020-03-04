package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/frontend"
	"github.com/buexplain/go-fool"
)

func frontend(mux *fool.Mux) {
	mux.Get("/", c_frontend.Index)
	mux.Get("/index-tag", c_frontend.IndexTag)
	mux.Get("/index-place", c_frontend.IndexPlace)
}
