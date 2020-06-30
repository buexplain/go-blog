package h_boot

import (
	"context"
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/cloudflare/tableflip"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func Run() {
	//初始化重启组件
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		a_boot.Logger.Error(err.Error())
		os.Exit(1)
	}
	defer upg.Stop()

	//监听信号
	go func() {
		sigCH := make(chan os.Signal, 1)
		signal.Notify(sigCH, []os.Signal{
			syscall.SIGHUP,
		}...)
		for range sigCH {
			err := upg.Upgrade()
			if err != nil {
				a_boot.Logger.ErrorF("upgrade failed: %w", err)
			}
		}
	}()

	//初始化tcp
	addr := a_boot.Config.App.Server.IP + ":" + strconv.Itoa(int(a_boot.Config.App.Server.Port))
	var ln net.Listener
	ln, err = upg.Listen("tcp", addr)
	if err != nil {
		a_boot.Logger.ErrorF("can't listen: %w", err)
		os.Exit(1)
	}

	//初始化http
	server := http.Server{
		Handler:      Server,
		ReadTimeout:  time.Duration(a_boot.Config.App.Server.ReadTimeout.Nanoseconds()),
		WriteTimeout: time.Duration(a_boot.Config.App.Server.WriteTimeout.Nanoseconds()),
		ErrorLog: NullLogger(),
	}
	//开启https后需要开启80端口做非https转发到https
	var redirectHttps *http.Server
	go func() {
		if a_boot.Config.App.Server.CertFile != "" && a_boot.Config.App.Server.KeyFile != "" {
			redirectHttps = NewRedirectHttps()
			go func() {
				<- time.After(time.Second*1)
				if err := redirectHttps.ListenAndServe(); err != nil && err != http.ErrServerClosed  {
					a_boot.Logger.Error(err.Error())
				}
			}()
			a_boot.Logger.Info("server running [pid " + strconv.Itoa(os.Getpid()) + "] " + "https://" + addr + "/backend/sign")
			if err := server.ServeTLS(ln, filepath.Join(a_boot.ROOT_PATH, a_boot.Config.App.Server.CertFile), filepath.Join(a_boot.ROOT_PATH, a_boot.Config.App.Server.KeyFile)); err != nil && err != http.ErrServerClosed {
				a_boot.Logger.Error(err.Error())
				os.Exit(1)
			}
		} else {
			a_boot.Logger.Info("server running [pid " + strconv.Itoa(os.Getpid()) + "] " + "http://" + addr + "/backend/sign")
			if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
				a_boot.Logger.Error(err.Error())
				os.Exit(1)
			}
		}
	}()

	//接受连接
	if err := upg.Ready(); err != nil {
		a_boot.Logger.Error(err.Error())
		os.Exit(1)
	}

	//等待结束
	<-upg.Exit()

	a_boot.Logger.Info("shutting down, please wait")

	//设定进程结束超时
	time.AfterFunc(time.Duration(a_boot.Config.App.Server.CloseTimedOut.Nanoseconds()), func() {
		a_boot.Logger.Error("graceful shutdown timed out")
		os.Exit(1)
	})

	//优先关闭跳转服务器
	if redirectHttps != nil {
		if err := redirectHttps.Shutdown(context.Background()); err != nil {
			a_boot.Logger.Error("shutdown redirect server err: " + err.Error())
		}
	}

	//关闭http
	if err := server.Shutdown(context.Background()); err != nil {
		a_boot.Logger.Error("graceful shutdown err: " + err.Error())
	}

	//等待事件调度器结束
	Bus.Close(time.Duration(a_boot.Config.App.Event.CloseTimedOut.Nanoseconds()))
	//等待日志收集器结束
	if err := Logger.Close(time.Duration(a_boot.Config.Log.CloseTimedOut.Nanoseconds())); err != nil {
		a_boot.Logger.Error(fmt.Errorf("h_boot.Logger.Close: %w", err).Error())
	}

	a_boot.Logger.Info("shutdown success")

	//退出程序
	os.Exit(0)
}
