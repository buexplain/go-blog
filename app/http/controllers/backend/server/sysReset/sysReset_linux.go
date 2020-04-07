package c_sysReset

import (
	"github.com/BurntSushi/toml"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/boot/config"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
	"io/ioutil"
	"os"
	"sync"
	"syscall"
	"time"
)

type lock struct {
	restarting int
	rw sync.RWMutex
}

var l *lock

func init()  {
	l = new(lock)
	l.restarting = -1
	l.rw = sync.RWMutex{}
}

func Start(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	//判断是否发送过信号
	if l.restarting > -1 {
		if l.restarting == 0 {
			return w.Error(1, "正在重启，请稍等...")
		}
		return w.Error(2, "重启失败，请登录服务器检查错误")
	}
	//接收客户端传递的配置
	c := r.Form("config")
	//检查配置格式是否正确
	if _, err := toml.Decode(c, &config.Config{}); err != nil {
		return w.Error(code.INVALID_ARGUMENT, code.Text(code.INVALID_ARGUMENT, err))
	}
	//写入配置到文件
	if err := ioutil.WriteFile(PATH, []byte(c), 0755); err != nil {
		return w.Error(code.SERVER, code.Text(code.SERVER, err))
	}
	//锁定
	l.restarting = 0
	go func() {
		<- time.After(a_boot.Config.App.Server.CloseTimedOut.Duration)
		l.restarting = 1
	}()
	//发送信号
	pid := os.Getpid()
	if err := syscall.Kill(pid, syscall.SIGHUP); err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	//返回进程id
	return w.Success(pid)
}