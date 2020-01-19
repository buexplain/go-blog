package s_services

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
	m_util "github.com/buexplain/go-blog/models/util"
	"strings"
)

//检查字段值是否存在
func CheckUnique(tableName string, field string, value interface{}, selfID ...int) bool {
	d := dao.Dao.Table("`" + tableName + "`")
	d.Where("`"+field+"`"+"=?", value)
	if len(selfID) > 0 && selfID[0] > 0 {
		d.Where("`ID`!=?", selfID[0])
	}
	type Tmp struct {
		ID                    int `xorm:"INTEGER"`
		models.DeletedAtField `xorm:"extends"`
	}
	has, err := d.Exist(new(Tmp))
	if err != nil {
		panic(err)
	}
	return !has
}

type Recursion struct {
	IdField string
	PidField string
	SortField string
	Select []string
	IsDown bool
}

func GetRecursion(table string, id int, rowsSlicePtr interface{}, recursion *Recursion) error {
	if recursion == nil {
		recursion = new(Recursion)
	}
	if recursion.IdField == "" {
		recursion.IdField = "ID"
	}
	if recursion.PidField == "" {
		recursion.PidField = "Pid"
	}
	if recursion.SortField == "" {
		recursion.SortField = "SortID"
	}
	if recursion.Select == nil || len(recursion.Select) == 0 {
		if recursion.Select == nil {
			recursion.Select = make([]string, 0)
		}
		t, err := m_util.GetTableInfo(dao.Dao, table)
		if err != nil {
			return err
		}else {
			for _, v := range t.Columns() {
				recursion.Select = append(recursion.Select, v.Name)
			}
		}
	}
	b := "`"+strings.Join(recursion.Select, "`,`")+"`"
	a := "`A`.`"+strings.Join(recursion.Select, "`,`A`.`")+"`"
	var sql string
	if recursion.IsDown {
		//向下查找所有的子级
		sql = `WITH tmp(%s)
		AS
		(
		SELECT %s FROM %s WHERE %s=%d
		UNION ALL
		SELECT %s FROM %s A, tmp B ON B.%s=A.%s ORDER BY %s ASC
		)
		select * from tmp`
		sql = fmt.Sprintf(sql, b, b, table, recursion.IdField, id, a, table, recursion.IdField, recursion.PidField, recursion.SortField)
	}else {
		//向上查找所有的父级
		sql = `WITH tmp(%s)
		AS
		(
		SELECT %s FROM %s WHERE %s=%d
		UNION ALL
		SELECT %s FROM %s A, tmp B ON A.%s=B.%s
		)
		select * from tmp ORDER BY %s ASC`
		sql = fmt.Sprintf(sql, b, b, table, recursion.IdField, id, a, table, recursion.IdField, recursion.PidField, recursion.IdField)
	}
	err := dao.Dao.NewSession().SQL(sql).Find(rowsSlicePtr)
	return err
}