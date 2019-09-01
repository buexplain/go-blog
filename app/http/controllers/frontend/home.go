package c_frontend

import (
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.Plain(http.StatusOK, "我与春风皆过客，你携秋水揽星河。")
}