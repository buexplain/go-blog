package c_sysReset

import (
	"github.com/buexplain/go-slim"
)

func Start(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	return w.Error(1, "当前系统不支持重启")
}
