package boot

import (
	"github.com/BurntSushi/toml"
	"github.com/buexplain/go-blog/app/boot/config"
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
		log.Panicln(err)
	}
	ROOT_PATH = filepath.ToSlash(filepath.Dir(dir))
	//兼容GoLand编辑器下的go run命令
	if strings.Contains(ROOT_PATH, "go-build") || strings.Contains(ROOT_PATH, "Temp") {
		ROOT_PATH = "./"
	}
}

//应用程序配置
var Config *config.Config

func init() {
	Config = new(config.Config)
	if _, err := toml.DecodeFile(filepath.Join(ROOT_PATH, "config.toml"), Config); err != nil {
		log.Panicln(err)
	}
}
