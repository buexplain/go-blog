package c_frontend

import (
	s_category "github.com/buexplain/go-blog/services/category"
	s_content "github.com/buexplain/go-blog/services/content"
	"github.com/buexplain/go-fool"
	"net/http"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	//categoryTree, err := s_category.GetTree()
	//if err != nil {
	//	return err
	//}
	//tagList, err := s_tag.GetALL()
	//if err != nil {
	//	return err
	//}
	//contentPlace, err := s_content.GetPlace()
	//if err != nil {
	//	return err
	//}

	page := r.QueryPositiveInt("page", 1)
	limit := r.QueryPositiveInt("limit", 15)
	category := r.QueryPositiveInt("category")
	place := r.Query("place")
	keyword := r.Query("keyword")
	contentList, counter, err := s_content.GetList(page, limit, category, place, keyword)
	if err != nil {
		return err
	}
	categoryNav := s_category.GetParents(category)

	return w.Success(contentList, counter, categoryNav)

	return w.Plain(http.StatusOK, "我与春风皆过客，你携秋水揽星河。")
}
