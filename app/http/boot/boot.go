package h_boot

import (
	"github.com/NYTimes/gziphandler"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/method"
	"github.com/buexplain/go-blog/app/http/boot/session"
	"github.com/buexplain/go-blog/app/http/boot/staticFile"
	"github.com/buexplain/go-blog/app/http/boot/trimForm"
	"github.com/buexplain/go-blog/app/http/boot/viewFunc"
	"github.com/buexplain/go-event"
	"github.com/buexplain/go-flog"
	"github.com/buexplain/go-flog/extra"
	"github.com/buexplain/go-flog/formatter"
	"github.com/buexplain/go-flog/handler"
	"github.com/buexplain/go-fool"
	"github.com/djherbis/fscache"
	"github.com/gorilla/csrf"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//设置日志
var Logger *flog.Logger

func init() {
	file := handler.NewFile(flog.GetLevelByName(a_boot.Config.Log.Level), formatter.NewLine(), filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Log.Path))
	//设置文件日志缓冲
	if a_boot.Config.Log.Buffer > 0 {
		file.SetBuffer(a_boot.Config.Log.Buffer, time.Duration(a_boot.Config.Log.Flush.Nanoseconds()))
	}
	Logger = flog.New(a_boot.Config.Log.Name, file).PushExtra(extra.NewFuncCaller())
	//debug模式下，打印日志到控制台
	if a_boot.Config.App.Debug {
		Logger.PushHandler(handler.NewSTD(flog.LEVEL_DEBUG, formatter.NewLine(), flog.LEVEL_WARNING))
	}
	if a_boot.Config.Log.Async {
		Logger.Async(int(a_boot.Config.Log.Capacity))
	}
}

//设置文件缓存
var Cache *fscache.FSCache

func init() {
	path := filepath.Join(a_boot.ROOT_PATH, a_boot.Config.Cache.Path)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalln(err)
	}
	if path == "." || path == "./" || path == `.\` {
		log.Fatalln("Root path not allowed")
	}
	c, err := fscache.New(path, 0755, time.Hour)
	if err != nil {
		log.Fatalln(err)
	}
	//先写一波文件，试试水
	err = ioutil.WriteFile(filepath.Join(path, "git.keep"), []byte(""), 0755)
	if err != nil {
		log.Fatalln(err)
	}
	Cache = c
}

//设置事件调度器
var Bus *event.Bus

func init() {
	Bus = event.New("http")
	if a_boot.Config.App.Event.Async {
		//开启事件异步处理
		Bus.Async(int(a_boot.Config.App.Event.Worker), int(a_boot.Config.App.Event.Capacity))
	}
}

//初始化app
var APP *fool.App

func init() {
	APP = fool.New(a_boot.Config.App.Debug)
}

//设置相关默认函数
func init() {
	//设置程序崩溃处理函数
	APP.SetRecoverFunc(defaultRecoverFunc)
	//设置错误处理函数
	APP.SetErrorFunc(defaultErrorFunc)
	//设置路由not found处理函数
	APP.Mux().SetDefaultRoute(defaultRoute)
}

//设置模板引擎
func init() {
	//设置模板路径
	APP.View().SetPath(filepath.Join(a_boot.ROOT_PATH, "resources/view"))
	//设置模板是否缓存
	APP.View().SetCache(!APP.Debug())
	//设置模板函数
	APP.View().AddFunc("message", viewFunc.Message)
	APP.View().AddFunc("URL", viewFunc.URL)
}

//设置session
func init() {
	switch a_boot.Config.Session.Store {
	case "cookie":
		h := session.NewCookieStoreManager(a_boot.Config.Session.Key)
		h.Options = a_boot.Config.Session.Options
		h.Name = a_boot.Config.Session.Name
		APP.SetSessionHandler(h)
		break
	case "file":
		path := filepath.Join(a_boot.ROOT_PATH, "storage/session")
		if err := os.MkdirAll(path, 0666); err != nil {
			log.Fatalln(err)
		}
		h := session.NewFilesystemStoreManager(path, a_boot.Config.Session.Key)
		h.Options = a_boot.Config.Session.Options
		h.Name = a_boot.Config.Session.Name
		APP.SetSessionHandler(h)
		break
	default:
		log.Fatalln("Config.Session.Store must set cookie|file")
	}
}

//设置全局中间件
func init() {
	//设置静态文件中间件
	if a_boot.Config.App.StaticFile.Enable {
		APP.Use(staticFile.Filter, http.MethodGet)
	}

	//设置方法欺骗中间件
	if a_boot.Config.App.Method.Enable {
		APP.Use(method.Filter, http.MethodGet, http.MethodPost)
	}

	//设置前后空白字符清理
	APP.Use(trimForm.Filter)
}

// http 服务器
var Server http.Handler

func init() {
	//设置csrf防御
	if a_boot.Config.CSRF.Enable {
		Server = csrf.Protect([]byte(a_boot.Config.CSRF.Key),
			csrf.ErrorHandler(defaultCSRFErrorHandler),
			csrf.CookieName(a_boot.Config.CSRF.Name),
			csrf.Path(a_boot.Config.CSRF.Options.Path),
			csrf.Domain(a_boot.Config.CSRF.Options.Domain),
			csrf.MaxAge(a_boot.Config.CSRF.Options.MaxAge),
			csrf.Secure(a_boot.Config.CSRF.Options.Secure),
			csrf.HttpOnly(a_boot.Config.CSRF.Options.HttpOnly),
			csrf.RequestHeader(a_boot.Config.CSRF.Header),
			csrf.FieldName(a_boot.Config.CSRF.Field))(APP)
	} else {
		Server = APP
	}
	//gzip压缩
	if a_boot.Config.GZIP.Enable {
		if wrap, err := gziphandler.NewGzipLevelHandler(a_boot.Config.GZIP.Level); err != nil {
			a_boot.Logger.ErrorF("gzip压缩的level配置错误: %s", err.Error())
			os.Exit(1)
		} else {
			Server = wrap(Server)
		}
	}
}
