package c_frontend

import (
	"github.com/buexplain/go-fool"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/http"
)


var engine *xorm.Engine
func init()  {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@tcp(192.168.158.40)/niuke_db?charset=utf8")
	if err != nil {
		panic(err)
	}
	engine.SetMaxOpenConns(40)
}

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	m, err := engine.Query("show tables")
	if err != nil {
		return w.Assign("code", 1).Assign("message", err.Error()).Assign("data", "").JSON(http.StatusOK)
	}
	return w.Assign("code", 0).Assign("message", "success").Assign("data", m).JSON(http.StatusOK)
}