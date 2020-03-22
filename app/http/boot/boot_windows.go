package h_boot

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
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
	server := http.Server{
		Addr:         addr,
		Handler:      Server,
		ReadTimeout:  time.Duration(a_boot.Config.App.Server.ReadTimeout.Nanoseconds()),
		WriteTimeout: time.Duration(a_boot.Config.App.Server.WriteTimeout.Nanoseconds()),
	}
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
