package c_official_user

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_user "github.com/buexplain/go-blog/models/user"
	m_util "github.com/buexplain/go-blog/models/util"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-validator"
	"github.com/gorilla/csrf"
	"net/http"
	"strconv"
)


//表单校验器
var v *validator.Validator

func init()  {
	v = validator.New()
	v.Rule("Account").Add("required", "请填写账号").Add("CheckUnique:id=0", "该账号已存在")
	v.Rule("Password",).Add("password:min=8,max=16",
		"请输入新密码",
		"新密码长度必须在8~16位之间",
		"密码格式有误，请输入数字、字母、符号",
		"密码格式有误，数字、字母、符号至少两种")
	v.Rule("Status").Add("in:in=1,2", "请选择状态")
	v.Rule("Identity").Add("in:in=1,2", "请选择身份")
	//校验账号是否存在
	v.Custom("CheckUnique", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, ok := value.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		if !m_util.CheckUnique("User", field, str, rule.GetInt("id")) {
			return rule.Message(0), nil
		}
		return "", nil
	})
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	query := m_util.NewQuery("User", ctx).Screen().Limit()
	query.Finder.Desc("ID")
	//设置查询条件后，先进行分页统计，//然后再进行连表查询，获取用户所有角色，避免跨表count
	count := query.Count()
	query.Finder.Join("LEFT", "`UserRoleRelation`", "`User`.`ID` = `UserRoleRelation`.`UserID`")
	query.Finder.Join("LEFT", "`Role`", "`UserRoleRelation`.`RoleID` = `Role`.`ID`")
	query.Finder.GroupBy("User.ID")
	query.Finder.Select("`User`.*, GROUP_CONCAT(`Role`.`Name`) as `RoleGroup`")

	type User struct {
		m_user.User `xorm:"extends"`
		RoleGroup string
	}
	var result []User
	query.Find(&result)

	if query.Error != nil {
		return query.Error
	}

	return w.
		Assign("count", count).
		Assign("result", result).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/user/index.html")
}

func Create(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	return w.
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").View(http.StatusOK, "backend/rbac/user/create.html")
}


func Store(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_user.User)

	if err := r.FormToStruct(mod); err != nil {
		return w.JumpBack(err)
	}

	if r, err := v.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if p, err := s_user.GeneratePassword(mod.Password); err != nil {
		return w.JumpBack(err)
	}else {
		mod.Password = p
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.JumpBack("操作成功")
}


func Edit(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := new(m_user.User)

	result.ID = r.ParamInt("id", 0)
	if result.ID <= 0 {
		return w.JumpBack("参数错误")
	}

	if ok, err := dao.Dao.Get(result); err != nil {
		return ctx.Error().WrapServer(err).Location()
	} else if !ok {
		return w.JumpBack("参数错误")
	}

	return w.
		Assign("result", result).
		Assign(boot.Config.CSRF.Field, csrf.TemplateField(r.Raw())).
		Layout("backend/layout/layout.html").
		View(http.StatusOK, "backend/rbac/user/create.html")
}


func Update(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	mod := new(m_user.User)
	if err := r.FormToStruct(mod); err != nil {
		return err
	}
	mod.ID = r.ParamInt("id", 0)

	vClone := v.Clone()
	vClone.Custom("password", func(field string, value interface{}, rule *validator.Rule, structVar interface{}) (s string, e error) {
		str, _ := value.(string)
		if str == "" {
			return "", nil
		}
		return validator.Pool("password")(field, value, rule, structVar)
	})
	vClone.Rule("ID").Add("required", "ID错误")
	vClone.Rule("Account").Add("CheckUnique:id="+strconv.Itoa(mod.ID), "该账号已存在")

	if r, err := vClone.Validate(mod); err != nil {
		return ctx.Error().WrapServer(err)
	}else if !r.IsEmpty() {
		return w.JumpBack(r)
	}

	if mod.Nickname == "" {
		mod.Nickname = mod.Account
	}

	if len(mod.Password) > 0 {
		if p, err := s_user.GeneratePassword(mod.Password); err != nil {
			return w.JumpBack(err)
		}else {
			mod.Password = p
		}
	}

	if _, err := dao.Dao.ID(mod.ID).Update(mod); err != nil {
		return ctx.Error().WrapServer(err).Location()
	}

	return w.Jump("/backend/rbac/user", "操作成功")
}