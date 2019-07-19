package config

type App struct {
	Debug      bool
	Server     Server
	Event      Event
	StaticFile StaticFile
	Method     Method
}

type Server struct {
	IP           string
	Port         uint
	ReadTimeout  Duration
	WriteTimeout Duration
}

type Event struct {
	Async    bool
	Worker   uint
	Capacity uint
}

type Method struct {
	Enable bool
	Field  string
}

type StaticFile struct {
	Enable  bool
	Path    []string
	Referer bool
}
