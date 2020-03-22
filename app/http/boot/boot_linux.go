package h_boot

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-gracehttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//启动http服务器
func Run() {
	addr := a_boot.Config.App.Server.IP + ":" + strconv.Itoa(int(a_boot.Config.App.Server.Port))
	Logger.Info("[pid " + strconv.Itoa(os.Getpid()) + "] " + "http://" + addr + "/backend/sign")
	server := gracehttp.NewServer(
		addr,
		Server,
		time.Duration(a_boot.Config.App.Server.ReadTimeout.Nanoseconds()),
		time.Duration(a_boot.Config.App.Server.WriteTimeout.Nanoseconds()),
	)
	server.SetErrorLogCallback(func(format string, args ...interface{}) {
		pid := strconv.Itoa(os.Getpid())
		format = "[pid " + pid + "] " + format
		Logger.Error(fmt.Sprintf(format, args...))
	})
	server.SetInfoLogCallback(func(format string, args ...interface{}) {
		pid := strconv.Itoa(os.Getpid())
		format = "[pid " + pid + "] " + format
		Logger.Info(fmt.Sprintf(format, args...))
	})
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Logger.Error(err.Error())
	}
	//等待事件调度器结束
	Bus.Close()
	//等待日志收集器结束
	if err := a_boot.Logger.Close(); err != nil {
		log.Println(err)
	}
	if err := Logger.Close(); err != nil {
		log.Println(err)
	}
}
