package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/frontend"
	"github.com/buexplain/go-blog/app/http/middleware"
	"github.com/buexplain/go-fool"
)

func frontend(mux *fool.Mux) {
	mux.Get("/", c_frontend.Index)
	mux.Get("/index.html", c_frontend.Index)
	mux.Get("/index-widget", c_frontend.IndexWidget).AddLabel("json")
	mux.Get("/article/:id.html", c_frontend.Article).Regexp("id.html", `^[1-9][0-9]*\.html$`).Use(middleware.CachePage)
	mux.Get("/article-hits/:id", c_frontend.ArticleHits)
}
