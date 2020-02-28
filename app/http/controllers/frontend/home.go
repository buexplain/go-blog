package c_frontend

import (
	"github.com/buexplain/go-blog/helpers"
	s_category "github.com/buexplain/go-blog/services/category"
	s_configItem "github.com/buexplain/go-blog/services/config/item"
	s_content "github.com/buexplain/go-blog/services/content"
	s_tag "github.com/buexplain/go-blog/services/tag"
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	//获取站点配置
	config := s_configItem.GetByGroup("SiteInfo")
	//获取站点菜单栏
	categoryTree := s_category.GetTree()
	//获取标签
	tagList := s_tag.GetALL()
	//获取归档
	contentPlace := s_content.GetPlace()
	//获取列表
	currentPage := r.QueryPositiveInt("page", 1)
	limit := r.QueryPositiveInt("limit", 10)
	keyword := r.Query("keyword")
	place := r.Query("place")
	tagID := r.QueryPositiveInt("tagID")
	categoryID := r.QueryPositiveInt("categoryID")
	contentList, total := s_content.GetList(currentPage, limit, categoryID, keyword,tagID, place)
	//注入数据
	w.Assign("config", config)
	w.Assign("categoryTree", categoryTree)
	w.Assign("tagList", tagList)
	w.Assign("contentPlace", contentPlace)
	w.Assign("contentList", contentList).Assign("total", total)
	w.Assign("pageHtml", helpers.PageHtmlSimple(*r.Raw().URL, currentPage, int(total), limit))
	w.Assign("tagID", tagID)
	w.Assign("keyword", keyword)
	w.Assign("place", place)
	//渲染模板
	return w.View(http.StatusOK, "frontend/index.html")
}
