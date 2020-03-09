package s_configGroup

import (
	"github.com/buexplain/go-blog/dao"
	m_configGroup "github.com/buexplain/go-blog/models/config/group"
	m_configItem "github.com/buexplain/go-blog/models/config/item"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-fool"
)

func GetList(ctx *fool.Ctx) (counter int64, result m_configGroup.List, err error) {
	query := s_services.NewQuery("ConfigGroup", ctx).Limit()
	query.Finder.Desc("ID")
	query.Where()
	//先获取分页所需的总条数
	counter = query.Count()
	query.Find(&result)
	err = query.Error
	return
}

func Destroy(ids []int) (int64, error) {
	affected, err := dao.Dao.In("ID", ids).Delete(new(m_configGroup.ConfigGroup))
	if err != nil {
		return 0, err
	}
	if affected > 0 {
		_, err = dao.Dao.In("GroupID", ids).Delete(new(m_configItem.ConfigItem))
		if err != nil {
			return 0, err
		}
	}
	return affected, nil
}
