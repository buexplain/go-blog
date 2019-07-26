package session

import (
	"encoding/base32"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
)

//session实例
type Session struct {
	s *sessions.Session
}

// 从session中读取一个条目
func (this *Session) Get(k interface{}) interface{} {
	return this.s.Values[k]
}

func (this *Session) GetString(k interface{}) string {
	if a, ok := this.s.Values[k].(string); ok {
		return a
	}
	return ""
}

func (this *Session) GetInt(k interface{}) int {
	if a, ok := this.s.Values[k].(int); ok {
		return a
	}
	return 0
}

func (this *Session) GetFloat64(k interface{}) float64 {
	if a, ok := this.s.Values[k].(float64); ok {
		return a
	}
	return 0
}

func (this *Session) GetFloat32(k interface{}) float32 {
	if a, ok := this.s.Values[k].(float32); ok {
		return a
	}
	return 0
}

// 从session中读取一个条目，并移除它
func (this *Session) Pull(k interface{}) interface{} {
	tmp := this.s.Values[k]
	delete(this.s.Values, k)
	return tmp
}

func (this *Session) PullString(k interface{}) string {
	a := this.s.Values[k]
	delete(this.s.Values, k)
	if a, ok := a.(string); ok {
		return a
	}
	return ""
}

func (this *Session) PullInt(k interface{}) int {
	a := this.s.Values[k]
	delete(this.s.Values, k)
	if a, ok := a.(int); ok {
		return a
	}
	return 0
}

func (this *Session) PullFloat64(k interface{}) float64 {
	a := this.s.Values[k]
	delete(this.s.Values, k)
	if a, ok := a.(float64); ok {
		return a
	}
	return 0
}

func (this *Session) PullFloat32(k interface{}) float32 {
	a := this.s.Values[k]
	delete(this.s.Values, k)
	if a, ok := a.(float32); ok {
		return a
	}
	return 0
}

// 增加一个session键值数据
func (this *Session) Set(k, v interface{}) {
	this.s.Values[k] = v
}

// 从session中移除一个条目
func (this *Session) Del(k interface{}) {
	delete(this.s.Values, k)
}

// 检查session里是否有此条目
func (this *Session) Has(k interface{}) bool {
	_, ok := this.s.Values[k]
	return ok
}

// 返回session id
func (this *Session) ID() string {
	return this.s.ID
}

//返回session name
func (this *Session) Name() string {
	return this.s.Name()
}

//重新生成一个session id
func (this *Session) Regenerate() {
	this.s.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
}

// 销毁session
func (this *Session) Destroy() {
	this.s.Options.MaxAge = -1
}

//session落地
func (this *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return this.s.Save(r, w)
}
