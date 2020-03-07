package dao

import (
	"github.com/buexplain/go-blog/app/boot"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

var Dao *xorm.Engine

func init() {
	db := filepath.Join(a_boot.ROOT_PATH, "database", a_boot.Config.Database.DSN)
	//创建目录
	if dir := filepath.Dir(db); dir != "" {
		if err := os.MkdirAll(dir, 777); err != nil {
			a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
			os.Exit(1)
		}
	}
	var err error
	Dao, err = NewDao(db)
	if err != nil {
		a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
	}
	//设置连接池的空闲数大小
	Dao.SetMaxIdleConns(10)
	//设置最大打开连接数
	Dao.SetMaxOpenConns(20)
	if a_boot.Config.App.Debug {
		Dao.Logger().SetLevel(core.LOG_INFO)
		Dao.ShowExecTime(true)
		Dao.ShowSQL(true)
	}
	//开启sqlite3的缓存
	_, err = Dao.Exec("PRAGMA cache_size = 5000")
	if err != nil {
		_ = Dao.Close()
		a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
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