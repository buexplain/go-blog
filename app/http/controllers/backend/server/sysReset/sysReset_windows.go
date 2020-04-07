package c_sysReset

import (
	"github.com/buexplain/go-fool"
)

func Start(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.Error(1, "当前系统不支持重启")
}
