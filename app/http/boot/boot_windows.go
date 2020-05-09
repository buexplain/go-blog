package h_boot

import (
	"context"
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

//启动http服务器
func Run() {
	//监听信号
	sigCH := make(chan os.Signal, 1)
	signal.Notify(sigCH, []os.Signal{
		syscall.SIGINT,
		syscall.SIGHUP,
	}...)

	//初始化tcp
	addr := a_boot.Config.App.Server.IP + ":" + strconv.Itoa(int(a_boot.Config.App.Server.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(fmt.Errorf("can't create new listener: %w", err))
	}

	//初始化http
	server := http.Server{
		Addr:         addr,
		Handler:      Server,
		ReadTimeout:  time.Duration(a_boot.Config.App.Server.ReadTimeout.Nanoseconds()),
		WriteTimeout: time.Duration(a_boot.Config.App.Server.WriteTimeout.Nanoseconds()),
	}
	go func() {
		if a_boot.Config.App.Server.CertFile != "" && a_boot.Config.App.Server.KeyFile != "" {
			log.Println("server running [pid " + strconv.Itoa(os.Getpid()) + "] " + "https://" + addr + "/backend/sign")
			if err := server.ServeTLS(ln, filepath.Join(a_boot.ROOT_PATH, a_boot.Config.App.Server.CertFile), filepath.Join(a_boot.ROOT_PATH, a_boot.Config.App.Server.KeyFile)); err != nil && err != http.ErrServerClosed {
				Logger.Error(err.Error())
			}
		}else {
			log.Println("server running [pid " + strconv.Itoa(os.Getpid()) + "] " + "http://" + addr + "/backend/sign")
			if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
				Logger.Error(err.Error())
			}
		}
	}()

	//等待信号
	for range sigCH {
		log.Println("shutting down, please wait")

		//设定进程结束超时
		time.AfterFunc(time.Duration(a_boot.Config.App.Server.CloseTimedOut.Nanoseconds()), func() {
			log.Fatalln("shutdown timed out")
		})

		//关闭http
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("shutdown err: "+err.Error())
		}

		//等待事件调度器结束
		Bus.Close(time.Duration(a_boot.Config.App.Event.CloseTimedOut.Nanoseconds()))
		//等待日志收集器结束
		if err := a_boot.Logger.Close(time.Duration(a_boot.Config.Log.CloseTimedOut.Nanoseconds())); err != nil {
			log.Println(fmt.Errorf("a_boot.Logger.Close: %w", err))
		}
		if err := Logger.Close(time.Duration(a_boot.Config.Log.CloseTimedOut.Nanoseconds())); err != nil {
			log.Println(fmt.Errorf("h_boot.Logger.Close: %w", err))
		}

		log.Println("shutdown success")

		//退出程序
		os.Exit(0)
	}
}
