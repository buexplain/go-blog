package c_frontend

import (
	"bytes"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_category "github.com/buexplain/go-blog/models/category"
	m_content "github.com/buexplain/go-blog/models/content"
	s_category "github.com/buexplain/go-blog/services/category"
	s_configItem "github.com/buexplain/go-blog/services/config/item"
	s_content "github.com/buexplain/go-blog/services/content"
	s_oauth "github.com/buexplain/go-blog/services/oauth"
	s_tag "github.com/buexplain/go-blog/services/tag"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-slim"
	"net/http"
	"strconv"
	"strings"
)

//首页
func Index(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	//获取站点配置
	config := s_configItem.GetByGroup("SiteInfo")
	//获取站点菜单栏
	categoryTree := s_category.GetTree(m_category.IsMenuYes)
	//获取列表
	categoryID := r.QueryPositiveInt("categoryID")
	tagID := r.QueryPositiveInt("tagID")
	place := r.Query("place")
	keyword := r.Query("keyword")
	currentPage := r.QueryPositiveInt("page", 1)
	limit := r.QueryPositiveInt("limit", 10)
	contentList := s_content.GetList(currentPage, limit, categoryID, tagID, place, keyword, m_content.OnlineYes)
	//注入数据
	w.Assign("config", config)
	w.Assign("categoryTree", categoryTree)
	if user := s_user.IsSignIn(r); user != nil {
		w.Assign("user", user)
	} else {
		w.Assign("github", s_oauth.NewGithub().GetURL("user", r.Raw().URL.String(), r))
	}
	w.Assign("contentList", contentList)
	w.Assign("limit", limit)
	w.Assign("prePage", currentPage-1)
	w.Assign("currentPage", currentPage)
	w.Assign("nextPage", currentPage+1)
	w.Assign("categoryID", categoryID)
	w.Assign("tagID", tagID)
	w.Assign("place", place)
	w.Assign("keyword", keyword)
	w.Assign("currentURL", *r.Raw().URL)
	//渲染模板
	return w.View(http.StatusOK, "frontend/index.html")
}

//首页的部件
func IndexWidget(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	buff := &bytes.Buffer{}
	result := map[string]string{}
	//渲染标签
	tagID := r.QueryPositiveInt("tagID")
	tagList := s_tag.GetALL()
	w.Assign("tagList", tagList)
	w.Assign("tagID", tagID)
	w.Assign("currentURL", *r.Raw().URL)
	buff.Reset()
	if err := w.Render(buff, w.Store().Pop(), "frontend/index-widget-tag.html"); err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	result["tag"] = buff.String()
	//渲染归档
	place := r.Query("place")
	contentPlace := s_content.GetPlace(m_content.OnlineYes)
	w.Assign("place", place)
	w.Assign("contentPlace", contentPlace)
	w.Assign("currentURL", *r.Raw().URL)
	buff.Reset()
	if err := w.Render(buff, w.Store().Pop(), "frontend/index-widget-place.html"); err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	result["place"] = buff.String()
	//返回结果
	return w.Success(result)
}

func Article(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	id, _ := strconv.Atoi(strings.TrimRight(r.Param("id.html"), ".html"))
	result, err := s_content.GetDetails(id, m_content.OnlineYes)
	if err != nil {
		return err
	}
	config := s_configItem.GetByGroup("SiteInfo")
	categoryTree := s_category.GetTree(0)
	w.Assign("currentURL", *r.Raw().URL)
	w.Assign("config", config)
	w.Assign("categoryTree", categoryTree)
	if user := s_user.IsSignIn(r); user != nil {
		w.Assign("user", user)
	} else {
		w.Assign("github", s_oauth.NewGithub().GetURL("user", r.Raw().URL.String(), r))
	}
	return w.Assign("result", result).View(http.StatusOK, "frontend/article.html")
}

func ArticleHits(ctx *slim.Ctx, w *slim.Response, r *slim.Request) error {
	contentID := r.ParamPositiveInt("id")
	isIncr := r.QueryBool("isIncr")
	if contentID == 0 {
		return w.Error(code.NOT_FOUND_DATA, code.Text(code.NOT_FOUND_DATA, contentID))
	}
	result := new(m_content.Content)
	has, err := dao.Dao.Where("ID=?", contentID).Select("Hits").Get(result)
	if err != nil {
		return w.Error(code.SERVER, err.Error())
	}
	if !has {
		return w.Error(code.NOT_FOUND_DATA, code.Text(code.NOT_FOUND_DATA, contentID))
	}
	if isIncr {
		result.Hits++
		_, err = dao.Dao.Table(result).ID(contentID).Update(map[string]interface{}{"Hits": result.Hits})
		if err != nil {
			return w.Error(code.SERVER, err.Error())
		}
	}
	return w.Success(result.Hits)
}
