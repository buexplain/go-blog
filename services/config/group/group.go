package s_configGroup

import (
	"github.com/buexplain/go-blog/dao"
	m_configGroup "github.com/buexplain/go-blog/models/config/group"
	m_configItem "github.com/buexplain/go-blog/models/config/item"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-slim"
)

func GetList(ctx *slim.Ctx) (counter int64, result m_configGroup.List, err error) {
	query := s_services.NewQuery("ConfigGroup", ctx)
	query.Finder.Desc("ID")
	query.Where().Limit().FindAndCount(&result, &counter)
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
