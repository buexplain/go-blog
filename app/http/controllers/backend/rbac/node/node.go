package c_node

import (
	"github.com/buexplain/go-blog/app/boot"
	e_syncRbacNode "github.com/buexplain/go-blog/app/http/events/syncRbacNode"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/node"
	"github.com/buexplain/go-blog/services/node"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

//表单校验器
var v *validator.Validator

func init()  {
	v = validator.New()
	v.Field("Name").Rule("required", "请输入节点名")
	v.Field("URL").Rule("required", "请输入访问路径")
	v.Field("Methods").Rule(`in:in=,GET,POST,PUT,DELETE&split=\,`, "请勾选请求方法", "错误的请求方法，请重新勾选")
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
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
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

	var restfulArr []*m_node.Node
	restful := r.FormBool("restful", false)
	if restful {
		restfulArr = make([]*m_node.Node, 0, 7)
		restfulArr = append(restfulArr, mod)

		create := new(m_node.Node)
		create.Name = "新增"
		create.URL = filepath.ToSlash(filepath.Join(mod.URL, "create"))
		create.SortID = 6
		create.IsMenu = m_node.IsMenuNo
		create.Methods = "GET"
		restfulArr = append(restfulArr, create)

		store := new(m_node.Node)
		store.Name = "保存"
		store.URL = mod.URL
		store.SortID = 5
		store.IsMenu = m_node.IsMenuNo
		store.Methods = "POST"
		restfulArr = append(restfulArr, store)

		edit := new(m_node.Node)
		edit.Name = "编辑"
		edit.URL = filepath.ToSlash(filepath.Join(mod.URL, "edit/:id"))
		edit.SortID = 4
		edit.IsMenu = m_node.IsMenuNo
		edit.Methods = "GET"
		restfulArr = append(restfulArr, edit)

		update := new(m_node.Node)
		update.Name = "更新"
		update.URL = filepath.ToSlash(filepath.Join(mod.URL, "update/:id"))
		update.SortID = 3
		update.IsMenu = m_node.IsMenuNo
		update.Methods = "PUT"
		restfulArr = append(restfulArr, update)

		destroy := new(m_node.Node)
		destroy.Name = "删除"
		destroy.URL = filepath.ToSlash(filepath.Join(mod.URL, "delete/:id"))
		destroy.SortID = 2
		destroy.IsMenu = m_node.IsMenuNo
		destroy.Methods = "DELETE"
		restfulArr = append(restfulArr, destroy)

		show := new(m_node.Node)
		show.Name = "查看"
		show.URL = filepath.ToSlash(filepath.Join(mod.URL, "show/:id"))
		show.SortID = 1
		show.IsMenu = m_node.IsMenuNo
		show.Methods = "GET"
		restfulArr = append(restfulArr, show)
	}else {
		restfulArr = make([]*m_node.Node, 0, 1)
		restfulArr = append(restfulArr, mod)
	}

	session := dao.Dao.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	for k, v := range restfulArr {
		if k != 0 {
			v.Pid = mod.ID
		}
		if _, err := session.Insert(v); err != nil {
			if err := session.Rollback(); err != nil {
				return ctx.Error().WrapServer(err).Location()
			}
			return ctx.Error().WrapServer(err)
		}
	}

	if err := session.Commit(); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}
	//触发超级角色的节点同步
	ctx.Event().Append(e_syncRbacNode.EVENT_NAME, a_boot.Config.Business.SuperRoleID)
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
		Assign(a_boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
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
	ids := r.QuerySliceInt("ids")
	if len(ids) == 0 {
		return w.JumpBack("参数错误")
	}
	if _, err := s_node.Destroy(ids); err != nil {
		return w.JumpBack(err)
	}
	//触发超级角色的节点同步
	ctx.Event().Append(e_syncRbacNode.EVENT_NAME, a_boot.Config.Business.SuperRoleID)
	return w.Jump("/backend/rbac/node", "操作成功")
}
