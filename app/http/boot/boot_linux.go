package boot

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
	addr := boot.Config.App.Server.IP + ":" + strconv.Itoa(int(boot.Config.App.Server.Port))
	Logger.Info("[pid " + strconv.Itoa(os.Getpid()) + "] " + "http://" + addr)
	server := gracehttp.NewServer(
		addr,
		Server,
		time.Duration(boot.Config.App.Server.ReadTimeout.Nanoseconds()),
		time.Duration(boot.Config.App.Server.WriteTimeout.Nanoseconds()),
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
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println(APP.Mux().DumpRouteMap())
	}()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Logger.Error(err.Error())
	}
	if err := APP.Close(); err != nil {
		log.Println(err)
	}
}
