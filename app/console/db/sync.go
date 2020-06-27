package db

import (
	"github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_attachment "github.com/buexplain/go-blog/models/attachment"
	m_category "github.com/buexplain/go-blog/models/category"
	m_configGroup "github.com/buexplain/go-blog/models/config/group"
	m_configItem "github.com/buexplain/go-blog/models/config/item"
	m_content "github.com/buexplain/go-blog/models/content"
	m_contentTag "github.com/buexplain/go-blog/models/contentTag"
	m_node "github.com/buexplain/go-blog/models/node"
	m_oauth "github.com/buexplain/go-blog/models/oauth"
	m_role "github.com/buexplain/go-blog/models/role"
	m_roleNodeRelation "github.com/buexplain/go-blog/models/roleNodeRelation"
	m_tag "github.com/buexplain/go-blog/models/tag"
	m_user "github.com/buexplain/go-blog/models/user"
	m_userRoleRelation "github.com/buexplain/go-blog/models/userRoleRelation"
	"github.com/spf13/cobra"
	"xorm.io/xorm"
)

//同步模型到数据库
var syncCmd *cobra.Command

func init() {
	syncCmd = &cobra.Command{
		Use:  "sync",
		Long: "同步models到数据库",
		Run: func(cmd *cobra.Command, args []string) {
			a_boot.Logger.Info("开始同步models到数据库")
			err := syncModels(dao.Dao)
			if err != nil {
				a_boot.Logger.ErrorF("同步models到数据库失败: %s", err)
			} else {
				a_boot.Logger.Info("同步models到数据库成功")
			}
		},
	}
	dbCmd.AddCommand(syncCmd)
}

func syncModels(engine *xorm.Engine) error {
	return engine.Sync2(
		new(m_node.Node),
		new(m_user.User),
		new(m_userRoleRelation.UserRoleRelation),
		new(m_role.Role),
		new(m_roleNodeRelation.RoleNodeRelation),
		new(m_category.Category),
		new(m_attachment.Attachment),
		new(m_tag.Tag),
		new(m_content.Content),
		new(m_contentTag.ContentTag),
		new(m_configGroup.ConfigGroup),
		new(m_configItem.ConfigItem),
		new(m_oauth.Oauth),
	)
}
