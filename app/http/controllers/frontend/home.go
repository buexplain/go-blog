package c_frontend

import (
	"bytes"
	"fmt"
	"github.com/buexplain/go-blog/app/http/boot/code"
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
	categoryID := r.QueryPositiveInt("categoryID")
	tagID := r.QueryPositiveInt("tagID")
	place := r.Query("place")
	keyword := r.Query("keyword")
	currentPage := r.QueryPositiveInt("page", 1)
	limit := r.QueryPositiveInt("limit", 10)
	debug.Set("列表")
	contentList := s_content.GetList(currentPage, limit, categoryID, tagID, place, keyword)
	debug.Get("列表")
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
	debug := NewDebug()
	debug.Set("标签")
	tagList := s_tag.GetALL()
	debug.Get("标签")
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
	contentPlace := s_content.GetPlace()
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
