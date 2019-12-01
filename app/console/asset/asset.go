package asset

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/app/console"
	"github.com/kevinburke/go-bindata"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var assetCmd *cobra.Command

func init() {
	assetCmd = &cobra.Command{
		Use:  "asset",
		Long: "静态资源打包、解包命令",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.Info("静态资源打包、解包命令")
		},
	}
	console.RootCmd.AddCommand(assetCmd)
}

//静态资源打包命令
var packCmd *cobra.Command

func init() {
	packCmd = &cobra.Command{
		Use: "pack",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.Info("开始打包静态资源")
			cfg := bindata.NewConfig()
			cfg.Package = "asset"
			cfg.Output = "app/console/asset/data.go"
			cfg.Input = make([]bindata.InputConfig, len(a_boot.Config.Asset.Dir))
			for i := range cfg.Input {
				cfg.Input[i] = parseInput(a_boot.Config.Asset.Dir[i])
			}
			err := bindata.Translate(cfg)
			if err != nil {
				a_boot.Logger.ErrorF("打包静态资源失败: %s", err)
				os.Exit(1)
			}
			a_boot.Logger.Info("打包静态资源成功")
		},
	}
	assetCmd.AddCommand(packCmd)
}

//静态资源解包命令
var unpackCmd *cobra.Command

func init() {
	unpackCmd = &cobra.Command{
		Use: "unpack",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.Info("开始解包静态资源")
			success := true
			for _, dir := range a_boot.Config.Asset.Dir {
				if strings.HasSuffix(dir, "/...") {
					dir = filepath.Clean(dir[:len(dir)-4])
				}
				if err := RestoreAssets(a_boot.ROOT_PATH, dir); err != nil {
					a_boot.Logger.Error(err.Error())
					success = false
				}
			}
			if success {
				a_boot.Logger.Info("解包静态资源成功")
			} else {
				a_boot.Logger.Error("解包静态资源失败")
				os.Exit(1)
			}
		},
	}
	assetCmd.AddCommand(unpackCmd)
}

// parseRecursive determines whether the given path has a recursive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//
//  ex:
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parseInput(path string) bindata.InputConfig {
	if strings.HasSuffix(path, "/...") {
		return bindata.InputConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	} else {
		return bindata.InputConfig{
			Path:      filepath.Clean(path),
			Recursive: false,
		}
	}
}
