package s_role

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
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
	delRoleSql := fmt.Sprintf("DELETE FROM Role WHERE %s AND ID NOT IN (SELECT Pid FROM Role WHERE %s)", in, notIn)

	in, err = builder.ToBoundSQL(builder.In("RoleID", ids))
	if err != nil {
		return 0, err
	}
	notIn, err = builder.ToBoundSQL(builder.In("ID", ids))
	if err != nil {
		return 0, err
	}
	delRoleNodeRelationSql := fmt.Sprintf("DELETE FROM RoleNodeRelation WHERE %s AND RoleID NOT IN ( SELECT ID FROM Role WHERE %s)", in, notIn)
	delUserRoleRelationSql := fmt.Sprintf("DELETE FROM UserRoleRelation WHERE %s AND RoleID NOT IN ( SELECT ID FROM Role WHERE %s)", in, notIn)

	//开启事务
	session := dao.Dao.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}
	//删除节点
	if result, err := session.Exec(delRoleSql); err != nil {
		if err := session.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}else {
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

			//删除用户角色表
			if _, err := session.Exec(delUserRoleRelationSql); err != nil {
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
