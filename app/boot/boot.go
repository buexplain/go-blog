package a_boot

import (
	"github.com/BurntSushi/toml"
	"github.com/buexplain/go-blog/app/boot/config"
	"github.com/buexplain/go-blog/helpers"
	"github.com/buexplain/go-flog"
	"github.com/buexplain/go-flog/extra"
	"github.com/buexplain/go-flog/formatter"
	"github.com/buexplain/go-flog/handler"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//程序根目录
var ROOT_PATH string

func init() {
	var dir string
	var err error
	//兼容GoLand编辑器下的go run命令
	tmp := strings.ToLower(os.Args[0])
	features := []string{"go_build", "go-build", "tmp", "temp"}
	isGoBuildRun := false
	for _, v := range features {
		if strings.Contains(tmp, v) {
			isGoBuildRun = true
			break
		}
	}
	if isGoBuildRun {
		dir, err = os.Getwd()
	} else {
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	if err != nil {
		log.Fatalln(err)
	}
	ROOT_PATH = strings.TrimSuffix(filepath.ToSlash(dir), "/") + "/"
}

//应用程序配置
var Config *config.Config

func init() {
	path := filepath.Join(ROOT_PATH, "config.toml")
	//读取配置文件
	c, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	//检查是否可以写入配置文件
	if err := ioutil.WriteFile(path, c, 0755); err != nil {
		log.Fatalln(err)
	}
	//解析配置文件到对象
	Config = new(config.Config)
	if _, err := toml.Decode(string(c), Config); err != nil {
		log.Fatalln(err)
	}
	//如果没有设置ip，则自动获取ip
	if Config.App.Server.IP == "" {
		Config.App.Server.IP = helpers.GetPublicIP()
	}
	//服务关闭的检查超时时间
	min := Config.Log.CloseTimedOut.Nanoseconds()*2 + Config.App.Event.CloseTimedOut.Nanoseconds()
	if Config.Log.Async {
		//日志开启了异步，要算上一次异步冲刷的时间
		min += Config.Log.Flush.Nanoseconds() * 2
	}
	//时间再翻倍
	min *= 2
	if Config.App.Server.CloseTimedOut.Nanoseconds() < min {
		tmp := time.Duration(min).String()
		if err := Config.App.Server.CloseTimedOut.UnmarshalText([]byte(tmp)); err != nil {
			log.Fatalln(err)
		}
	}
	//如果证书不为空，强制设置端口为443
	if Config.App.Server.CertFile != "" && Config.App.Server.KeyFile != "" {
		Config.App.Server.Port = 443
	}
}

//全局共用的控制台日志
//这个日志无缓冲
var Logger *flog.Logger

func init() {
	Logger = flog.New("std", handler.NewSTD(flog.LEVEL_DEBUG, formatter.NewLine().SetTimeFormat("2006-01-02 15:04:05.99"), flog.LEVEL_ERROR))
	Logger.PushExtra(extra.NewFuncCaller())
}
