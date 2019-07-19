package dao

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"time"
	"xorm.io/core"
)

var Dao *xorm.Engine

func init() {
	db := filepath.Join(boot.ROOT_PATH, "database/database.db")
	var err error
	Dao, err = NewDao(db)
	if err != nil {
		boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
	}
	if boot.Config.App.Debug {
		Dao.ShowSQL(true)
		Dao.Logger().SetLevel(core.LOG_DEBUG)
	}
}

func NewDao(path string) (*xorm.Engine, error) {
	dao, err := xorm.NewEngine("sqlite3", path)
	if err != nil {
		return nil, err
	}
	//设置时区
	dao.TZLocation = time.Local
	dao.DatabaseTZ = time.Local
	//设置结构体与表字段一致
	dao.SetMapper(core.SameMapper{})
	return dao, nil
}
