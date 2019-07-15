package session

import (
	"errors"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"path/filepath"
)

var Options sessions.Options

func init() {
	Options = sessions.Options{
		Path:     "",
		Domain:   "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
}

type Manager struct {
	CookieStore     *sessions.CookieStore
	FilesystemStore *sessions.FilesystemStore
	Name            string
	Options         sessions.Options
}

func (this *Manager) Get(r *fool.Request) (fool.Session, error) {
	var s *sessions.Session
	var err error
	if this.CookieStore != nil {
		s, err = this.CookieStore.Get(r.Raw(), this.Name)
	} else if this.FilesystemStore != nil {
		s, err = this.FilesystemStore.Get(r.Raw(), this.Name)
	} else {
		err = errors.New("not found session store")
	}
	if s != nil {
		var o sessions.Options
		o = this.Options
		s.Options = &o
		return &Session{s: s}, nil
	}
	return nil, err
}

func NewCookieStoreManager(key string) *Manager {
	tmp := new(Manager)
	tmp.Name = "session"
	tmp.CookieStore = sessions.NewCookieStore([]byte(key))
	tmp.Options = Options
	return tmp
}

func NewFilesystemStoreManager(path string, key string) *Manager {
	for ; len(key) < 16; key += key {
	}
	tmp := new(Manager)
	tmp.Name = "session"
	tmp.FilesystemStore = sessions.NewFilesystemStore(filepath.Join(path))
	tmp.FilesystemStore.Codecs = []securecookie.Codec{securecookie.New([]byte(key), []byte(key)[0:16])}
	tmp.Options = Options
	return tmp
}
