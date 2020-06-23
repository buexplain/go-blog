package config

//日志配置
type Log struct {
	Async         bool
	Capacity      uint
	Name          string
	Path          string
	Level         string
	Buffer        int
	Flush         Duration
	CloseTimedOut Duration
}
