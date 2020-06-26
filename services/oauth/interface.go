package s_oauth

import (
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	"github.com/buexplain/go-slim"
)

type AccessResult interface {
	GetAccessToken() string
}

type UserInfo interface {
	GetAccount() string
	GetNickname() string
}

type Oauth interface {
	GetURL(scope string, redirect string, r *slim.Request) (auth_url string)
	GetAccessToken(r *slim.Request) (AccessResult, error)
	GetUserInfo(access_token string) (UserInfo, error)
	GetStatus() m_oauth.Status
}
