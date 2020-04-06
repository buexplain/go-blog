package c_sysRestart

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
	"os"
	"syscall"
)

func Restart(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if err := syscall.Kill(os.Getgid(), syscall.SIGHUP); err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	return w.Success("重启成功")
}