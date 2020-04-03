package c_sysInfo

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-gracehttp"
	"os"
	"syscall"
)

func Restart(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	if os.Getenv(gracehttp.GRACEFUL_ENVIRON_KEY) != "" {
		return w.Success("重启成功")
	}
	err := syscall.Kill(os.Getgid(), syscall.SIGUSR2)
	if err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	return w.Success("重启成功")
}