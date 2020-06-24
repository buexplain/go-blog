package s_configItem

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	m_configItem "github.com/buexplain/go-blog/models/config/item"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-fool"
	"html/template"
)

func GetList(ctx *fool.Ctx) (counter int64, result m_configItem.List, err error) {
	query := s_services.NewQuery("ConfigItem", ctx)
	query.Finder.Desc("ID")
	groupID := ctx.Request().QueryPositiveInt("groupID")
	if groupID > 0 {
		query.Finder.Where("GroupID=?", groupID)
	}
	query.Where().Limit().FindAndCount(&result, &counter)
	err = query.Error
	return
}

func Destroy(ids []int) (int64, error) {
	affected, err := dao.Dao.In("ID", ids).Delete(new(m_configItem.ConfigItem))
	if err != nil {
		return 0, err
	}
	return affected, nil
}

type List struct {
	group string
	data  m_configItem.List
}

func (this List) Get(key string) string {
	for _, v := range this.data {
		if v.Key == key {
			return v.Value
		}
	}
	panic(code.NewF(code.INVALID_CONFIG, "请在管理后台设置: %s%s%s", this.group,".", key))
}

func (this List) GetToHTML(key string) template.HTML {
	for _, v := range this.data {
		if v.Key == key {
			return template.HTML(v.Value)
		}
	}
	panic(code.NewF(code.INVALID_CONFIG, "请在管理后台设置: %s%s%s", this.group,".", key))
}

func GetByGroup(groupName string) *List {
	mod := dao.Dao.Table("ConfigItem")
	mod.Select("ConfigItem.`Key`, ConfigItem.`Value`")
	mod.Join("inner", "ConfigGroup", "ConfigGroup.ID = ConfigItem.GroupID")
	mod.Where("ConfigGroup.`Key`=?", groupName)
	result := &m_configItem.List{}
	err := mod.Find(result)
	if err != nil {
		panic(err)
	}
	if len(*result) == 0 {
		panic(code.NewF(code.INVALID_CONFIG, "请在管理后台设置: %s", groupName))
	}
	return &List{data: *result, group: groupName}
}
