package s_tag

import (
	"github.com/buexplain/go-blog/dao"
	m_tag "github.com/buexplain/go-blog/models/tag"
)

func Store(name string) (int, error) {
	mod := &m_tag.Tag{}
	mod.Name = name
	b, err := dao.Dao.Get(mod)
	if err != nil {
		return 0, err
	}
	if b {
		return mod.ID, nil
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return 0, err
	}

	return mod.ID, nil
}