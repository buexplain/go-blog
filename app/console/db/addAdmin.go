package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_user "github.com/buexplain/go-blog/models/user"
	m_util "github.com/buexplain/go-blog/models/util"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/spf13/cobra"
	"os"
)

//添加一个管理员用户
var addAdminCmd *cobra.Command

func init() {
	var account string
	var password string
	addAdminCmd = &cobra.Command{
		Use:  "addAdmin",
		Long: "新增管理员",
		Run: func(cmd *cobra.Command, args []string) {
			boot.Logger.InfoF("开始新增管理员: %s %s", account, password)
			if len(account) == 0 {
				boot.Logger.Error("缺失参数: --account")
				os.Exit(1)
			}
			if len(password) == 0 {
				boot.Logger.Error("缺失参数: --password")
				os.Exit(1)
			}

			user := m_user.User{}
			user.Account = account
			var err error
			user.Password, err = s_user.GeneratePassword(password)
			if err != nil {
				boot.Logger.ErrorF("新增用户失败: %s", err.Error())
			}
			user.Status = m_user.StatusAllow
			user.Identity = m_user.IdentityOfficial
			user.Nickname = account

			if !m_util.CheckUnique("User", "Account", "account") {
				boot.Logger.ErrorF("该用户已存在: %s", account)
				os.Exit(1)
			}

			if _, err := dao.Dao.Insert(user); err != nil {
				boot.Logger.ErrorF("新增用户失败: %s", err.Error())
				os.Exit(1)
			}

			boot.Logger.Info("新增管理员成功")
		},
	}
	addAdminCmd.Flags().StringVarP(&account, "account", "a", "", "账号")
	addAdminCmd.Flags().StringVarP(&password, "password", "p", "", "密码")
	dbCmd.AddCommand(addAdminCmd)
}