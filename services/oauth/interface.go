package s_oauth

import (
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	"github.com/buexplain/go-fool"
)

type AccessResult interface {
	GetAccessToken() string
}

type UserInfo interface {
	GetAccount() string
	GetNickname() string
}

type Oauth interface {
	GetURL(scope string, redirect string, r *fool.Request) (auth_url string)
	GetAccessToken(r *fool.Request) (AccessResult, error)
	GetUserInfo(access_token string) (UserInfo, error)
	GetStatus() m_oauth.Status
}
