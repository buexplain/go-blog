package config

type Log struct {
	Async bool
	Capacity uint
	Name string
	Path string
	Level string
	Buffer int
	Flush Duration
}