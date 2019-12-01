package h_boot

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/method"
	"github.com/buexplain/go-blog/app/http/boot/session"
	"github.com/buexplain/go-blog/app/http/boot/staticFile"
	"github.com/buexplain/go-blog/app/http/boot/viewFunc"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/flog"
	"github.com/buexplain/go-fool/flog/extra"
	"github.com/buexplain/go-fool/flog/formatter"
	"github.com/buexplain/go-fool/flog/handler"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var APP *fool.App

//初始化app
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
	APP.SetLoggerHandler(Logger)
}

//设置模板引擎
func init() {
	//设置模板路径
	APP.View().SetPath(filepath.Join(a_boot.ROOT_PATH, "resources/view"))
	//设置模板是否缓存
	APP.View().SetCache(!APP.Debug())
	//设置模板函数
	APP.View().AddFunc("message", viewFunc.Message)
}

//设置事件监听者
func init() {
	if a_boot.Config.App.Event.Async {
		//开启事件异步处理
		APP.EventHandler().Async(int(a_boot.Config.App.Event.Worker), int(a_boot.Config.App.Event.Capacity))
	}
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
		h := session.NewFilesystemStoreManager(filepath.Join(a_boot.ROOT_PATH, "storage/session"), a_boot.Config.Session.Key)
		h.Options = a_boot.Config.Session.Options
		h.Name = a_boot.Config.Session.Name
		APP.SetSessionHandler(h)
		break
	default:
		log.Panicln("Config.Session.Store must set cookie|file")
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
}

// http 服务器
var Server http.Handler

func init() {
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
}
