package c_frontend

import (
	"fmt"
	"github.com/buexplain/go-blog/helpers"
	s_category "github.com/buexplain/go-blog/services/category"
	s_configItem "github.com/buexplain/go-blog/services/config/item"
	s_content "github.com/buexplain/go-blog/services/content"
	s_tag "github.com/buexplain/go-blog/services/tag"
	"github.com/buexplain/go-fool"
	"net/http"
	"time"
)

type Debug struct {
	data map[string]time.Time
}

func NewDebug() *Debug  {
	return &Debug{data: map[string]time.Time{}}
}

func (this *Debug) Set(name string) {
	this.data[name] = time.Now()
}

func (this Debug) Get(name string) {
	if t, ok := this.data[name]; ok {
		fmt.Printf("%s 耗时: %d 毫秒\n", name, (time.Now().UnixNano() - t.UnixNano())/1e6)
	}
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	debug := NewDebug()
	//获取站点配置
	debug.Set("配置")
	config := s_configItem.GetByGroup("SiteInfo")
	debug.Get("配置")
	//获取站点菜单栏
	debug.Set("导航")
	categoryTree := s_category.GetTree()
	debug.Get("导航")
	//获取列表
	currentPage := r.QueryPositiveInt("page", 1)
	limit := r.QueryPositiveInt("limit", 10)
	keyword := r.Query("keyword")
	place := r.Query("place")
	tagID := r.QueryPositiveInt("tagID")
	categoryID := r.QueryPositiveInt("categoryID")
	debug.Set("列表")
	contentList := s_content.GetList(currentPage, limit, categoryID, keyword, tagID, place)
	debug.Get("列表")
	//注入数据
	w.Assign("config", config)
	w.Assign("categoryTree", categoryTree)
	w.Assign("contentList", contentList)
	w.Assign("pageHtml", helpers.PageHtmlSimple(*r.Raw().URL, currentPage, len(contentList), limit))
	w.Assign("tagID", tagID)
	w.Assign("keyword", keyword)
	w.Assign("place", place)
	//渲染模板
	return w.View(http.StatusOK, "frontend/index.html")
}

func IndexTag(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	tagID := r.QueryPositiveInt("tagID")
	debug := NewDebug()
	debug.Set("标签")
	tagList := s_tag.GetALL()
	debug.Get("标签")
	w.Assign("tagList", tagList)
	w.Assign("tagID", tagID)
	return w.View(http.StatusOK, "frontend/index-tag.html")
}

func IndexPlace(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	place := r.Query("place")
	debug := NewDebug()
	debug.Set("归档")
	contentPlace := s_content.GetPlace()
	debug.Get("归档")
	w.Assign("place", place)
	w.Assign("contentPlace", contentPlace)
	return w.View(http.StatusOK, "frontend/index-place.html")
}