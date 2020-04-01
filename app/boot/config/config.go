package config

type Config struct {
	App      App
	Log      Log
	Cache    Cache
	Session  Session
	CSRF     CSRF
	GZIP     GZIP
	Asset    Asset
	Business Business
	Database Database
}
