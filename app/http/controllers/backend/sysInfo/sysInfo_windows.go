package c_sysInfo

import "github.com/buexplain/go-fool"

func Restart(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.Error(1, "当前操作系统不支持平滑重启")
}