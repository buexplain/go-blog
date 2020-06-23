package s_oauth

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	"github.com/buexplain/go-fool"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Github struct {
	ID          string
	Secret      string
	RedirectUri string
}

func NewGithub() Oauth {
	tmp := &Github{}
	if config, ok := a_boot.Config.Business.OAuth.List["github"]; ok {
		tmp.ID = config.ID
		tmp.Secret = config.Secret
		tmp.RedirectUri = setQueryString(config.CallBackUrl, "third_site", string(ThirdSiteGithub))
	}
	return tmp
}

//@link https://developer.github.com/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/
func (this Github) GetURL(scope string, redirect string, r *fool.Request) (oauth_url string) {
	if this.ID == "" || this.Secret == "" {
		return
	}
	tmp := url.URL{}
	tmp.Host = "github.com"
	tmp.Scheme = "https"
	tmp.Path = "login/oauth/authorize"
	query := url.Values{}
	query.Set("client_id", this.ID)
	query.Set("redirect_uri", this.RedirectUri)
	query.Set("scope", scope)
	state := randomString(6)
	//记录state
	r.Session().Set("oauthState", state)
	//记录授权成功后的跳转地址
	r.Session().Set("oauthRedirect", redirect)
	//记录授权失败后的跳转地址
	r.Session().Set("oauthOrigin", r.Raw().URL.String())
	query.Set("state", state)
	tmp.RawQuery = query.Encode()
	oauth_url = tmp.String()
	return
}

type GithubAccessResult struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

func (this GithubAccessResult) GetAccessToken() string {
	return this.AccessToken
}

func (this Github) GetAccessToken(r *fool.Request) (AccessResult, error) {
	state := r.Query("state")
	codeStr := r.Query("code")
	if state == "" {
		return nil, code.NewM(code.INVALID_ARGUMENT, "state")
	}
	if codeStr == "" {
		return nil, code.NewM(code.INVALID_ARGUMENT, "code")
	}
	//校验state
	if r.Session().GetString("oauthState") != state {
		return nil, code.NewM(code.INVALID_ARGUMENT, "state")
	}
	if this.ID == "" {
		return nil, code.NewM(code.INVALID_CONFIG, "github id")
	}
	if this.Secret == "" {
		return nil, code.NewM(code.INVALID_CONFIG, "github secret")
	}
	form := url.Values{}
	form.Set("client_id", this.ID)
	form.Set("client_secret", this.Secret)
	form.Set("code", codeStr)
	form.Set("state", state)
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, code.New(code.CALL_THIRD_ERROR, "请求github错误，请重试")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", "name=anny")
	client := &http.Client{}
	client.Timeout = time.Second * 10
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result GithubAccessResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type GithubUser struct {
	AvatarURL         string      `json:"avatar_url,omitempty"`
	Bio               interface{} `json:"bio,omitempty"`
	Blog              string      `json:"blog,omitempty"`
	Collaborators     int64       `json:"collaborators,omitempty"`
	Company           interface{} `json:"company,omitempty"`
	CreatedAt         string      `json:"created_at,omitempty"`
	DiskUsage         int64       `json:"disk_usage,omitempty"`
	Email             interface{} `json:"email,omitempty"`
	EventsURL         string      `json:"events_url,omitempty"`
	Followers         int64       `json:"followers,omitempty"`
	FollowersURL      string      `json:"followers_url,omitempty"`
	Following         int64       `json:"following,omitempty"`
	FollowingURL      string      `json:"following_url,omitempty"`
	GistsURL          string      `json:"gists_url,omitempty"`
	GravatarID        string      `json:"gravatar_id,omitempty"`
	Hireable          interface{} `json:"hireable,omitempty"`
	HTMLURL           string      `json:"html_url,omitempty"`
	ID                int64       `json:"id,omitempty"`
	Location          interface{} `json:"location,omitempty"`
	Login             string      `json:"login,omitempty"`
	Name              interface{} `json:"name,omitempty"`
	NodeID            string      `json:"node_id,omitempty"`
	OrganizationsURL  string      `json:"organizations_url,omitempty"`
	OwnedPrivateRepos int64       `json:"owned_private_repos,omitempty"`
	Plan              struct {
		Collaborators int64  `json:"collaborators,omitempty"`
		Name          string `json:"name,omitempty"`
		PrivateRepos  int64  `json:"private_repos,omitempty"`
		Space         int64  `json:"space,omitempty"`
	} `json:"plan,omitempty"`
	PrivateGists            int64  `json:"private_gists,omitempty"`
	PublicGists             int64  `json:"public_gists,omitempty"`
	PublicRepos             int64  `json:"public_repos,omitempty"`
	ReceivedEventsURL       string `json:"received_events_url,omitempty"`
	ReposURL                string `json:"repos_url,omitempty"`
	SiteAdmin               bool   `json:"site_admin,omitempty"`
	StarredURL              string `json:"starred_url,omitempty"`
	SubscriptionsURL        string `json:"subscriptions_url,omitempty"`
	TotalPrivateRepos       int64  `json:"total_private_repos,omitempty"`
	TwoFactorAuthentication bool   `json:"two_factor_authentication,omitempty"`
	Type                    string `json:"type,omitempty"`
	UpdatedAt               string `json:"updated_at,omitempty"`
	URL                     string `json:"url,omitempty"`
}

func (this GithubUser) GetAccount() string {
	return this.Login
}

func (this GithubUser) GetNickname() string {
	return this.Login
}

func init() {
	gob.Register(&GithubUser{})
}

func (this Github) GetUserInfo(access_token string) (UserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/user?access_token=%s", access_token))
	if err != nil {
		return nil, code.New(code.CALL_THIRD_ERROR, "请求github错误，请重试")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result GithubUser
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (this Github) GetStatus() m_oauth.Status {
	return m_oauth.StatusGithub
}
