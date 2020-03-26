package c_frontend

import (
	"bytes"
	"github.com/buexplain/go-blog/app/http/boot/code"
	m_category "github.com/buexplain/go-blog/models/category"
	m_content "github.com/buexplain/go-blog/models/content"
	s_category "github.com/buexplain/go-blog/services/category"
	s_configItem "github.com/buexplain/go-blog/services/config/item"
	s_content "github.com/buexplain/go-blog/services/content"
	s_tag "github.com/buexplain/go-blog/services/tag"
	"github.com/buexplain/go-fool"
	"net/http"
	"strconv"
	"strings"
)

//首页
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
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
func IndexWidget(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
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

func Article(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
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
	return w.Assign("result", result).View(http.StatusOK, "frontend/article.html")
}
