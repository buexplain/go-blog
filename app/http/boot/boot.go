package boot

import (
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/flog"
	"github.com/buexplain/go-fool/flog/extra"
	"github.com/buexplain/go-fool/flog/formatter"
	"github.com/buexplain/go-fool/flog/handler"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/method"
	"github.com/buexplain/go-blog/app/http/boot/session"
	"github.com/buexplain/go-blog/app/http/boot/staticFile"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var APP *fool.App

//初始化app
func init() {
	APP = fool.New(boot.Config.App.Debug)
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
	file := handler.NewFile(flog.GetLevelByName(boot.Config.Log.Level), formatter.NewLine(), filepath.Join(boot.ROOT_PATH, boot.Config.Log.Path))
	//设置文件日志缓冲
	if boot.Config.Log.Buffer > 0 {
		file.SetBuffer(boot.Config.Log.Buffer, time.Duration(boot.Config.Log.Flush.Nanoseconds()))
	}
	Logger = flog.New(boot.Config.Log.Name, file).PushExtra(extra.NewFuncCaller())
	//debug模式下，打印日志到控制台
	if boot.Config.App.Debug {
		Logger.PushHandler(handler.NewSTD(flog.LEVEL_DEBUG, formatter.NewLine(), flog.LEVEL_WARNING))
	}
	if boot.Config.Log.Async {
		Logger.Async(int(boot.Config.Log.Capacity))
	}
	APP.SetLoggerHandler(Logger)
}

//设置模板引擎
func init() {
	//设置模板路径
	APP.View().SetPath(filepath.Join(boot.ROOT_PATH, "resources/view"))
	//设置模板是否缓存
	APP.View().SetCache(!APP.Debug())
}

//设置事件监听者
func init() {
	if boot.Config.App.Event.Async {
		//开启事件异步处理
		APP.EventHandler().Async(int(boot.Config.App.Event.Worker), int(boot.Config.App.Event.Capacity))
	}
}

//设置session
func init() {
	switch boot.Config.Session.Store {
	case "cookie":
		h := session.NewCookieStoreManager(boot.Config.Session.Key)
		h.Options = boot.Config.Session.Options
		h.Name = boot.Config.Session.Name
		APP.SetSessionHandler(h)
		break
	case "file":
		h := session.NewFilesystemStoreManager(filepath.Join(boot.ROOT_PATH, "storage/session"), boot.Config.Session.Key)
		h.Options = boot.Config.Session.Options
		h.Name = boot.Config.Session.Name
		APP.SetSessionHandler(h)
		break
	default:
		log.Panicln("Config.Session.Store must set cookie|file")
	}
}

//设置全局中间件
func init() {
	//设置静态文件中间件
	if boot.Config.App.StaticFile.Enable {
		APP.Use(staticFile.Filter, http.MethodGet)
	}

	//设置方法模拟中间件
	if boot.Config.App.Method.Enable {
		APP.Use(method.Filter, http.MethodGet, http.MethodPut, http.MethodPatch)
	}
}

//设置全局路由正则
func init() {
	APP.Mux().Regexp("id", `^[1-9][0-9]*$`)
}

// http 服务器
var Server http.Handler
func init()  {
	if boot.Config.CSRF.Enable {
		Server = csrf.Protect([]byte(boot.Config.CSRF.Key),
			csrf.ErrorHandler(defaultCSRFErrorHandler),
			csrf.CookieName(boot.Config.CSRF.Name),
			csrf.Path(boot.Config.CSRF.Options.Path),
			csrf.Domain(boot.Config.CSRF.Options.Domain),
			csrf.MaxAge(boot.Config.CSRF.Options.MaxAge),
			csrf.Secure(boot.Config.CSRF.Options.Secure),
			csrf.HttpOnly(boot.Config.CSRF.Options.HttpOnly),
			csrf.RequestHeader(boot.Config.CSRF.Header),
			csrf.FieldName(boot.Config.CSRF.Field))(APP)
	}else {
		Server = APP
	}
}