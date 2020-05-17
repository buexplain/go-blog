package s_oauth

import (
	"fmt"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
)

type ThirdSite string

const ThirdSiteGithub ThirdSite = "github"

func New(r *fool.Request) (Oauth, error) {
	third_site := ThirdSite(r.Query("third_site"))
	switch third_site {
	case ThirdSiteGithub:
		return NewGithub(), nil
	default:
		return nil, errors.MarkClient(fmt.Errorf("invalid request param: third_site"))
	}
}