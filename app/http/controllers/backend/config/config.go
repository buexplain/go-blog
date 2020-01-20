package c_config

import (
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.View(http.StatusOK, "backend/config/index.html")
}

func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.View(http.StatusOK, "backend/config/index.html")
}