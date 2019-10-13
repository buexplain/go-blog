package config

//csrf防御配置
type CSRF struct {
	Enable  bool
	Key     string
	Name    string
	Field   string
	Header  string
	Options struct {
		Path     string
		Domain   string
		MaxAge   int
		Secure   bool
		HttpOnly bool
	}
}
