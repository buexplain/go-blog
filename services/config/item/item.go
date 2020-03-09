package s_configItem

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	m_configItem "github.com/buexplain/go-blog/models/config/item"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"html/template"
)

func GetList(ctx *fool.Ctx) (counter int64, result m_configItem.List, err error) {
	query := s_services.NewQuery("ConfigItem", ctx).Limit()
	query.Finder.Desc("ID")
	query.Where()
	groupID := ctx.Request().QueryPositiveInt("groupID")
	if groupID > 0 {
		query.Finder.Where("GroupID=?", groupID)
	}
	//先获取分页所需的总条数
	counter = query.Count()
	query.Find(&result)
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
	panic(errors.MarkServer(fmt.Errorf("not found config: %s.%s", this.group, key)))
}

func (this List) GetToHTML(key string) template.HTML {
	for _, v := range this.data {
		if v.Key == key {
			return template.HTML(v.Value)
		}
	}
	panic(errors.MarkServer(fmt.Errorf("not found config: %s.%s", this.group, key)))
}

func GetByGroup(groupName string) *List {
	//SELECT ConfigItem.* FROM ConfigItem INNER JOIN ConfigGroup on ConfigGroup.ID = ConfigItem.GroupID WHERE ConfigGroup."Key" = "SiteInfo"
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
		panic(errors.MarkServer(fmt.Errorf("not found config: %s", groupName)))
	}
	return &List{data: *result, group: groupName}
}
