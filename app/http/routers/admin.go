package routers

import (
	"github.com/buexplain/go-blog/app/http/controllers/admin/sign"
	"github.com/buexplain/go-fool"
)

func admin(mux *fool.Mux) {
	mux.Get("/admin/sign", c_sign.Index)
	mux.Post("/admin/sign", c_sign.In)
}