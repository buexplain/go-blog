package c_sysRestart

import (
	"github.com/buexplain/go-fool"
	"net/http"
)

func Restart(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.View(http.StatusOK, "backend/server/sysRestart/index.html")
}
