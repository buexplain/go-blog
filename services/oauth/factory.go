package s_oauth

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-fool"
)

type ThirdSite string

const ThirdSiteGithub ThirdSite = "github"

func New(r *fool.Request) (Oauth, error) {
	third_site := ThirdSite(r.Query("third_site"))
	switch third_site {
	case ThirdSiteGithub:
		return NewGithub(), nil
	default:
		return nil, code.NewM(code.INVALID_ARGUMENT, "third_site")
	}
}
