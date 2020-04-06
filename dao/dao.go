package dao

import (
	"fmt"
	"github.com/buexplain/go-blog/app/boot"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
	"os"
	"path/filepath"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var Dao *xorm.Engine

func forceDSNParam(dsn string, param map[string]string) (string, error) {
	s, err := url.PathUnescape(dsn)
	if err != nil {
		return "", err
	}
	if u, err := url.Parse(s); err != nil {
		return "", err
	} else {
		q := u.Query()
		for k, v := range param {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		return u.String(), nil
	}
}

func init() {
	//强制sqlite3的时区为本地时区
	dsn, err := forceDSNParam(a_boot.Config.Database.DSN, map[string]string{"_loc": "auto"})
	if err != nil {
		a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
	}

	//创建目录
	db := filepath.Join(a_boot.ROOT_PATH, "database", dsn)
	if dir := filepath.Dir(db); dir != "" {
		if err := os.MkdirAll(dir, 777); err != nil {
			a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
			os.Exit(1)
		}
	}

	//打开数据库
	Dao, err = xorm.NewEngine("sqlite3", db)
	if err != nil {
		a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
	}

	//强制xorm按UTC时区存储时间
	Dao.DatabaseTZ = time.UTC

	//强制xorm的程序时区设置为本地时区
	Dao.TZLocation = time.Local

	//设置结构体与表字段一致
	Dao.SetMapper(core.SameMapper{})

	//设置连接池中的最大闲置连接数
	if a_boot.Config.Database.MaxIdleConns > 0 {
		Dao.SetMaxIdleConns(a_boot.Config.Database.MaxIdleConns)
	}

	//设置与数据库建立连接的最大数目
	if a_boot.Config.Database.MaxOpenConns > 0 {
		Dao.SetMaxOpenConns(a_boot.Config.Database.MaxOpenConns)
	}

	//是否开启日志打印
	if a_boot.Config.Database.ShowSQL {
		Dao.Logger().SetLevel(log.LOG_DEBUG)
		Dao.ShowSQL(true)
	}

	//开启sqlite3的缓存
	_, err = Dao.Exec(fmt.Sprintf("PRAGMA cache_size = %+d", a_boot.Config.Database.CacheSize))
	if err != nil {
		_ = Dao.Close()
		a_boot.Logger.ErrorF("初始化dao失败: %s", err.Error())
		os.Exit(1)
	}
}
