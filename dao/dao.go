package dao

import (
	"github.com/buexplain/go-blog/app/boot"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var Dao *xorm.Engine

func init() {
	if strings.Index(a_boot.Config.Database.DSN, "_loc") != -1 {
		a_boot.Logger.ErrorF("初始化dao失败: 请不要设置参数`_loc`, 数据库默认按UTC时区存储时间")
		os.Exit(1)
	}
	db := filepath.Join(a_boot.ROOT_PATH, "database", a_boot.Config.Database.DSN)
	//创建目录
	if dir := filepath.Dir(db); dir != "" {
		if err := os.MkdirAll(dir, 777); err != nil {
			a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
			os.Exit(1)
		}
	}

	//打开数据库
	var err error
	Dao, err = xorm.NewEngine("sqlite3", db)
	if err != nil {
		a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
	}

	//默认为UTC时区存储时间
	Dao.DatabaseTZ = time.UTC

	//设置结构体与表字段一致
	Dao.SetMapper(core.SameMapper{})

	//设置连接池的空闲数大小
	Dao.SetMaxIdleConns(10)

	//设置最大打开连接数
	Dao.SetMaxOpenConns(20)

	if a_boot.Config.App.Debug {
		Dao.Logger().SetLevel(log.LOG_DEBUG)
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