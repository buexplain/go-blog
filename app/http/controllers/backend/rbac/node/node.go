package c_node

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/node"
	"github.com/buexplain/go-blog/services/node"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"strings"
)

//表单校验器
var v *validator.Validator

func init()  {
	v = validator.New()
	v.Rule("Name").Add("required", "请输入节点名")
	v.Rule("URL").Add("required", "请输入访问路径")
	v.Rule("Methods").Add(`in:in=,GET,POST,PUT,DELETE&split=\,`, "请勾选请求方法", "错误的请求方法，请重新勾选")
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
	mod := &m_node.Node{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}
	mod.Methods = strings.Join(r.FormSlice("methods", make([]string, 0)), ",")

	if r, err := v.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

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
	mod := &m_node.Node{}
	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	mod.ID = r.ParamInt("id", 0)
	if mod.ID <= 0 {
		return w.JumpBack("参数错误")
	}
	mod.Methods = strings.Join(r.FormSlice("methods", make([]string, 0)), ",")

	if r, err := v.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

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
