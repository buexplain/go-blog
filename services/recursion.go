package s_services

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"strings"
	"xorm.io/builder"
)

//获取父级
func GetParents(table string, id int, rowsSlicePtr interface{}, field ...string) error {
	if len(field) == 0 {
		t, err := GetTableInfo(dao.Dao, table)
		if err != nil {
			return err
		}else {
			for _, v := range t.Columns() {
				field = append(field, v.Name)
			}
		}
	}
	b := "`"+strings.Join(field, "`,`")+"`"
	a := "`A`.`"+strings.Join(field, "`,`A`.`")+"`"
	//向上查找所有的父级
	sql := `WITH tmp(%s)
		AS
		(
		SELECT %s FROM %s WHERE ID=%d
		UNION ALL
		SELECT %s FROM %s A, tmp B ON A.ID=B.Pid
		)
		select * from tmp ORDER BY ID ASC`
	sql = fmt.Sprintf(sql, b, b, table, id, a, table)
	return dao.Dao.NewSession().SQL(sql).Find(rowsSlicePtr)
}

//获取子级
func GetSons(table string, id int, rowsSlicePtr interface{}, sortField string, where builder.Cond, field ...string) error {
	if sortField == "" {
		sortField = "ID"
	}
	if len(field) == 0 {
		t, err := GetTableInfo(dao.Dao, table)
		if err != nil {
			return err
		}else {
			for _, v := range t.Columns() {
				field = append(field, v.Name)
			}
		}
	}
	b := "`"+strings.Join(field, "`,`")+"`"
	a := "`A`.`"+strings.Join(field, "`,`A`.`")+"`"
	//向下查找所有的子级
	sql := `WITH tmp(%s)
		AS
		(
		SELECT %s FROM %s WHERE ID=%d
		UNION ALL
		SELECT %s FROM %s A, tmp B ON B.ID=A.Pid ORDER BY %s ASC
		)
		select * from tmp`
	sql = fmt.Sprintf(sql, b, b, table, id, a, table, sortField)
	if where != nil {
		if s, err := builder.ToBoundSQL(where); err != nil {
			return err
		}else {
			sql += " where " + s
		}
	}
	return dao.Dao.NewSession().SQL(sql).Find(rowsSlicePtr)
}