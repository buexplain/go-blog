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
	IP       string
	Port     uint
	CertFile string
	KeyFile  string
	//http请求的表单编码类型为multipart/form-data的内容解析到内存中的大小，超出会解析到磁盘
	FormMaxMemory int64
	//http请求的body的大小限制
	BodyMaxBytes  int64
	ReadTimeout   Duration
	WriteTimeout  Duration
	CloseTimedOut Duration
}

//事件调度配置
type Event struct {
	Async         bool
	Worker        uint
	Capacity      uint
	CloseTimedOut Duration
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
