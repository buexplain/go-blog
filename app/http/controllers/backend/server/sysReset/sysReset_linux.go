package c_sysReset

import (
	"github.com/BurntSushi/toml"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/boot/config"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-slim"
	"io/ioutil"
	"os"
	"sync"
	"syscall"
	"time"
)

type lock struct {
	restarting int
	rw         sync.RWMutex
}

var l *lock

func init() {
	l = new(lock)
	l.restarting = -1
	l.rw = sync.RWMutex{}
}

func Start(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	//判断是否发送过信号
	if l.restarting > -1 {
		if l.restarting == 0 {
			return w.Error(1, "正在重启，请稍等...")
		}
		return w.Error(2, "重启失败，请登录服务器检查错误日志")
	}
	//接收客户端传递的配置
	c := r.Form("config")
	//检查配置格式是否正确
	if _, err := toml.Decode(c, &config.Config{}); err != nil {
		return w.Error(code.INVALID_ARGUMENT, err.Error())
	}
	//写入配置到文件
	if err := ioutil.WriteFile(PATH, []byte(c), 0755); err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	//锁定
	l.restarting = 0
	go func() {
		//进程退出超时的时候还没结束的话，锁定重启按钮，不再支持重启
		<-time.After(a_boot.Config.App.Server.CloseTimedOut.Duration)
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
