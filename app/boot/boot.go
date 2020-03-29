package a_boot

import (
	"github.com/BurntSushi/toml"
	"github.com/buexplain/go-blog/app/boot/config"
	"github.com/buexplain/go-blog/helpers"
	"github.com/buexplain/go-flog"
	"github.com/buexplain/go-flog/extra"
	"github.com/buexplain/go-flog/formatter"
	"github.com/buexplain/go-flog/handler"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//程序根目录
var ROOT_PATH string

func init() {
	dir, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatalln(err)
	}
	ROOT_PATH = filepath.ToSlash(filepath.Dir(dir))
	//兼容GoLand编辑器下的go run命令
	if strings.Contains(ROOT_PATH, "go-build") || strings.Contains(ROOT_PATH, "Temp") {
		ROOT_PATH = "./"
	}else {
		ROOT_PATH = strings.TrimSuffix(ROOT_PATH, "/")+"/"
	}
}

//应用程序配置
var Config *config.Config

func init() {
	Config = new(config.Config)
	if _, err := toml.DecodeFile(filepath.Join(ROOT_PATH, "config.toml"), Config); err != nil {
		log.Fatalln(err)
	}
	if Config.App.Server.IP == "" {
		Config.App.Server.IP = helpers.GetPublicIP()
	}
}

//全局共用的控制台日志
var Logger *flog.Logger

func init() {
	Logger = flog.New("std", handler.NewSTD(flog.LEVEL_DEBUG, formatter.NewLine().SetTimeFormat("2006-01-02 15:04:05.99"), flog.LEVEL_ERROR))
	Logger.PushExtra(extra.NewFuncCaller())
}
