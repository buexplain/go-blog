package c_node

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/node"
	s_node "github.com/buexplain/go-blog/services/node"
	"github.com/buexplain/go-fool"
	"github.com/gorilla/csrf"
	"github.com/thedevsaddam/govalidator"
	"html/template"
	"net/http"
	"strings"
)

func init()  {
	govalidator.AddCustomRule("methods", func(field string, rule string, message string, value interface{}) error {
		fmt.Println(field)
		fmt.Println(rule)
		fmt.Println(message)
		fmt.Println(value)
		return nil
	})
}

//列表
func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result, err := m_node.GetALL()
	if err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	return w.
		Assign("result", template.JS(result.String())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/node/index.html")
}

//创建
func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	pid := r.ParamInt("pid", 0)
	return ctx.Response().
		Assign("pid", pid).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/node/create.html")
}

//保存
func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	rules := govalidator.MapData{
		"Name": []string{"required"},
		"URL": []string{"required"},
		"Methods": []string{"methods"},
	}

	messages := govalidator.MapData{
		"Name": []string{"required:请输入节点名"},
		"URL": []string{"required:请输入访问路径"},
		"Methods": []string{"methods:请勾选请求方法"},
	}

	opts := govalidator.Options{
		Request:         r.Raw(),
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) > 0 {
		return w.JumpBack(e)
	}

	mod := &m_node.Node{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	mod.Methods = strings.Join(r.FormSlice("methods", make([]string, 0)), ",")


	if _, err := mod.Insert(); err != nil {
		return w.JumpBack(err)
	}

	return w.JumpBack("操作成功")
}

//编辑
func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_node.Node)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if has, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err)
	} else if !has {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/node/create.html")
}

//更新
func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	rules := govalidator.MapData{
		"Name": []string{"required"},
		"URL": []string{"required"},
		"Methods": []string{"required"},
	}

	messages := govalidator.MapData{
		"Name": []string{"required:请输入节点名"},
		"URL": []string{"required:请输入访问路径"},
		"Methods": []string{"required:请勾选请求方法"},
	}

	opts := govalidator.Options{
		Request:         r.Raw(),
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) > 0 {
		return w.JumpBack(e)
	}

	mod := &m_node.Node{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	mod.ID = r.ParamInt("id", 0)
	if mod.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	mod.Methods = strings.Join(r.FormSlice("methods", make([]string, 0)), ",")

	if _, err := mod.Update(); err != nil {
		return w.JumpBack(err)
	}

	return w.Jump("/backend/rbac/node", "操作成功")
}

//删除
func Destroy(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	ids := r.QuerySliceInt("ids[]")
	if len(ids) == 0 {
		return w.JumpBack("参数错误")
	}
	if _, err := s_node.Destroy(ids); err != nil {
		return w.JumpBack(err)
	}
	return w.Jump("/backend/rbac/node", "操作成功")
}
