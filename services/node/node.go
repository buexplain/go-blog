package s_node

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	m_node "github.com/buexplain/go-blog/models/node"
	"strings"
	"xorm.io/builder"
)

func Destroy(ids []int) (int64, error) {
	var err error
	var in string
	var notIn string
	in, err = builder.ToBoundSQL(builder.In("ID", ids))
	if err != nil {
		return 0, err
	}
	notIn, err = builder.ToBoundSQL(builder.In("Pid", ids))
	if err != nil {
		return 0, err
	}
	delNodeSql := fmt.Sprintf("DELETE FROM Node WHERE %s AND ID NOT IN (SELECT Pid FROM Node WHERE %s)", in, notIn)

	in, err = builder.ToBoundSQL(builder.In("NodeID", ids))
	if err != nil {
		return 0, err
	}
	notIn, err = builder.ToBoundSQL(builder.In("ID", ids))
	if err != nil {
		return 0, err
	}
	delRoleNodeRelationSql := fmt.Sprintf("DELETE FROM RoleNodeRelation WHERE %s AND NodeID NOT IN ( SELECT ID FROM Node WHERE %s)", in, notIn)

	//开启事务
	session := dao.Dao.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}
	//删除节点
	if result, err := session.Exec(delNodeSql); err != nil {
		if err := session.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	} else {
		//获取受影响的行数
		var affected int64
		if affected, err = result.RowsAffected(); err != nil {
			if err := session.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
		if affected > 0 {
			//删除角色节点表
			if _, err := session.Exec(delRoleNodeRelationSql); err != nil {
				if err := session.Rollback(); err != nil {
					return 0, err
				}
				return 0, err
			}
		}
		//提交事务
		if err := session.Commit(); err != nil {
			return 0, err
		}
		//返回删除结果
		return affected, nil
	}
}

func Store(mod *m_node.Node) (affected int64, err error) {
	if mod.IsMenu == 0 {
		mod.IsMenu = m_node.IsMenuNo
	} else {
		mod.IsMenu = m_node.IsMenuYes
	}
	mod.URL = strings.Trim(mod.URL, " ")
	if mod.ID == 0 {
		return dao.Dao.MustCols("Pid", "Methods").Insert(mod)
	}
	return dao.Dao.ID(mod.ID).MustCols("Pid", "Methods").Update(mod)
}

//根据url获取node
func GetByURL(path string) *m_node.Node {
	result := new(m_node.Node)
	has, err := dao.Dao.Table("Node").Where("URL=?", path).Get(result)
	if err != nil || !has {
		return nil
	}
	return result
}