package s_oauth

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Github struct {
	ID string
	Secret string
	RedirectUri string
}

func NewGithub() Oauth {
	tmp := &Github{}
	if config, ok := a_boot.Config.Business.OAuth["github"]; ok {
		tmp.ID = config.ID
		tmp.Secret = config.Secret
		tmp.RedirectUri = setQueryString(config.RedirectUri, "third_site", string(ThirdSiteGithub))
	}
	return tmp
}

//@link https://developer.github.com/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/
func (this Github) GetURL(scope string, oauth_after_url string, r *fool.Request) (oauth_url string) {
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
	r.Session().Set("oauth_state", state)
	r.Session().Set("oauth_after_url", oauth_after_url)
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
	code := r.Query("code")
	if state == "" {
		return nil, errors.MarkClient(fmt.Errorf("invalid request param: state"))
	}
	if code == "" {
		return nil, errors.MarkClient(fmt.Errorf("invalid request param: code"))
	}
	//校验state
	if r.Session().GetString("oauth_state") != state {
		return nil, errors.MarkClient(fmt.Errorf("invalid request param: state"))
	}
	if this.ID == "" || this.Secret == "" {
		return nil, fmt.Errorf("not found config: github")
	}
	form := url.Values{}
	form.Set("client_id", this.ID)
	form.Set("client_secret", this.Secret)
	form.Set("code", code)
	form.Set("state", state)
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", "name=anny")
	client := &http.Client{}
	client.Timeout = time.Second*3
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
	if err := json.Unmarshal(body,  &result); err != nil {
		return nil, err
	}
	return &result, nil
}


type GithubUser struct {
	Login                   string    `json:"login,omitempty"`
	ID                      int       `json:"id,omitempty"`
	NodeID                  string    `json:"node_id,omitempty"`
	AvatarURL               string    `json:"avatar_url,omitempty"`
	GravatarID              string    `json:"gravatar_id,omitempty"`
	URL                     string    `json:"url,omitempty"`
	HtmlURL                 string    `json:"html_url,omitempty"`
	FollowersURL            string    `json:"followers_url,omitempty"`
	FollowingURL            string    `json:"following_url,omitempty"`
	GistsURL                string    `json:"gists_url,omitempty"`
	StarredURL              string    `json:"starred_url,omitempty"`
	SubscriptionsURL        string    `json:"subscriptions_url,omitempty"`
	OrganizationsURL        string    `json:"organizations_url,omitempty"`
	ReposURL                string    `json:"repos_url,omitempty"`
	EventsURL               string    `json:"events_url,omitempty"`
	ReceivedEventsURL       string    `json:"received_events_url,omitempty"`
	Type                    string    `json:"type,omitempty"`
	SiteAdmin               bool      `json:"site_admin,omitempty"`
	Name                    string    `json:"name,omitempty"`
	Company                 string    `json:"company,omitempty"`
	Blog                    string    `json:"blog,omitempty"`
	Location                string    `json:"location,omitempty"`
	Email                   string    `json:"email,omitempty"`
	Hireable                string    `json:"hireable,omitempty"`
	Bio                     string    `json:"bio,omitempty"`
	PublicRepos             string    `json:"public_repos,omitempty"`
	PublicGists             int       `json:"public_gists,omitempty"`
	Followers               int       `json:"followers,omitempty"`
	Following               int       `json:"following,omitempty"`
	CreatedAt               time.Time `json:"created_at,omitempty"`
	UpdatedAt               time.Time `json:"updated_at,omitempty"`
	PrivateGists            int       `json:"private_gists,omitempty"`
	TotalPrivateRepos       int       `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos       int       `json:"owned_private_repos,omitempty"`
	DiskUsage               int       `json:"disk_usage,omitempty"`
	Collaborators           int       `json:"collabocollaboratorsrators,omitempty"`
	TwoFactorAuthentication bool      `json:"two_factor_authentication,omitempty"`
}

func (this GithubUser) GetAccount() string {
	return this.Name
}

func (this GithubUser) GetNickname() string {
	return this.Name
}

func init() {
	gob.Register(&GithubUser{})
}

func (this Github) GetUserInfo(access_token string) (UserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/user?access_token=%s", access_token))
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
	var result GithubUser
	h_boot.Logger.Info(string(body))
	if err := json.Unmarshal(body,  &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (this Github) GetStatus() m_oauth.Status {
	return m_oauth.StatusGithub
}