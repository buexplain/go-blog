package boot

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
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
	server := http.Server{
		Addr:         addr,
		Handler:      Server,
		ReadTimeout:  time.Duration(boot.Config.App.Server.ReadTimeout.Nanoseconds()),
		WriteTimeout: time.Duration(boot.Config.App.Server.WriteTimeout.Nanoseconds()),
	}
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println(APP.Mux().DumpRouteMap("github.com/buexplain/go-blog/"))
	}()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Logger.Error(err.Error())
	}
	if err := APP.Close(); err != nil {
		log.Println(err)
	}
}
