package config

type App struct {
	Debug      bool
	Server     Server
	Event      Event
	StaticFile StaticFile
	Method     Method
}

//http服务配置
type Server struct {
	IP           string
	Port         uint
	ReadTimeout  Duration
	WriteTimeout Duration
}

//事件调度配置
type Event struct {
	Async    bool
	Worker   uint
	Capacity uint
}

//restful方法欺骗中间件配置
type Method struct {
	Enable bool
	Field  string
}

//静态文件配置
type StaticFile struct {
	Enable  bool
	Path    []string
	Referer bool
}
