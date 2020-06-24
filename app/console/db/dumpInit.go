package db

import (
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/dao"
	m_models "github.com/buexplain/go-blog/models"
	m_configGroup "github.com/buexplain/go-blog/models/config/group"
	m_configItem "github.com/buexplain/go-blog/models/config/item"
	m_node "github.com/buexplain/go-blog/models/node"
	m_role "github.com/buexplain/go-blog/models/role"
	m_roleNodeRelation "github.com/buexplain/go-blog/models/roleNodeRelation"
	m_user "github.com/buexplain/go-blog/models/user"
	m_userRoleRelation "github.com/buexplain/go-blog/models/userRoleRelation"
	s_services "github.com/buexplain/go-blog/services"
	s_user "github.com/buexplain/go-blog/services/user"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

//导出初始化的init.sql
var dumpInitCmd *cobra.Command

func init() {
	dumpInitCmd = &cobra.Command{
		Use:  "dumpInit",
		Long: "导出数据库到 ./database/init.sql 文件",
		Run: func(cmd *cobra.Command, args []string) {
			//打开内存数据库
			memoryDao, err := xorm.NewEngine("sqlite3", ":memory:")
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				defer func() {
					_ = memoryDao.Close()
				}()
				//强制xorm按UTC时区存储时间
				memoryDao.DatabaseTZ = time.UTC
				//强制xorm的程序时区设置为本地时区
				memoryDao.TZLocation = time.Local
				//设置结构体与表字段一致
				memoryDao.SetMapper(core.SameMapper{})
			}
			//同步models到数据库
			if err := syncModels(memoryDao); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			//同步节点
			var nodeList m_node.List
			nodeList, err = m_node.GetALL()
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				for _, v := range nodeList {
					if _, err := memoryDao.Insert(v); err != nil {
						a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
						os.Exit(1)
					}
				}
			}
			//插入角色
			adminRole := &m_role.Role{}
			adminRole.Name = "超级管理员"
			adminRole.SortID = 1991
			adminRole.Pid = 0
			if _, err := memoryDao.Insert(adminRole); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				roleNodeRelationList := make(m_roleNodeRelation.List, 0, len(nodeList))
				for _,v := range nodeList {
					roleNodeRelationList = append(roleNodeRelationList, m_roleNodeRelation.RoleNodeRelation{RoleID:adminRole.ID, NodeID:v.ID})
				}
				if _, err := memoryDao.Insert(roleNodeRelationList); err != nil {
					a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
					os.Exit(1)
				}
			}
			guestRole := &m_role.Role{}
			guestRole.Name = "来宾"
			guestRole.SortID = 1991
			guestRole.Pid = 0
			if _, err := memoryDao.Insert(guestRole); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				roleNodeRelationList := make(m_roleNodeRelation.List, 0, len(nodeList))
				denyNodeList := map[string]bool{
					"/backend/backup/start":true,
					"/backend/backup/download":true,
					"/backend/backup/delete":true,
				}
				for _, v := range nodeList {
					//跳过非 GET 的 路由
					if !strings.EqualFold(v.Methods, http.MethodGet) {
						continue
					}
					//跳过不允许的路由
					if _, ok := denyNodeList[v.URL]; ok {
						continue
					}
					roleNodeRelationList = append(roleNodeRelationList, m_roleNodeRelation.RoleNodeRelation{RoleID:guestRole.ID, NodeID:v.ID})
				}
				if _, err := memoryDao.Insert(roleNodeRelationList); err != nil {
					a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
					os.Exit(1)
				}
			}
			//插入用户表
			adminUser := &m_user.User{}
			adminUser.Account = "admin"
			adminUser.Nickname = "admin"
			adminUser.Password, err = s_user.GeneratePassword("123456")
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			adminUser.Identity = m_user.IdentityOfficial
			adminUser.Status = m_user.StatusAllow
			adminUser.LastTime = m_models.Time(time.Now())
			if _, err := memoryDao.Insert(adminUser); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				userRoleRelation := m_userRoleRelation.UserRoleRelation{}
				userRoleRelation.UserID = adminUser.ID
				userRoleRelation.RoleID = adminRole.ID
				if _, err := memoryDao.Insert(userRoleRelation); err != nil {
					a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
					os.Exit(1)
				}
			}
			guestUser := &m_user.User{}
			guestUser.Account = "guest"
			guestUser.Nickname = "guest"
			guestUser.Password, err = s_user.GeneratePassword("123456")
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			guestUser.Identity = m_user.IdentityOfficial
			guestUser.Status = m_user.StatusAllow
			guestUser.LastTime = m_models.Time(time.Now())
			if _, err := memoryDao.Insert(guestUser); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				userRoleRelation := m_userRoleRelation.UserRoleRelation{}
				userRoleRelation.UserID = guestUser.ID
				userRoleRelation.RoleID = guestRole.ID
				if _, err := memoryDao.Insert(userRoleRelation); err != nil {
					a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
					os.Exit(1)
				}
			}
			//插入站点配置
			siteInfo := &m_configGroup.ConfigGroup{}
			siteInfo.Name = "站点信息"
			siteInfo.Key = "SiteInfo"
			siteInfo.Comment = "本站点基本信息"
			if _, err := memoryDao.Insert(siteInfo); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}else {
				siteInfoItem := m_configItem.List{}
				siteInfoItem = append(siteInfoItem, &m_configItem.ConfigItem{GroupID:siteInfo.ID, Name:"站点名称", Key:"Name", Value:"梦想星辰大海"})
				siteInfoItem = append(siteInfoItem, &m_configItem.ConfigItem{GroupID:siteInfo.ID, Name:"站点关键词", Key:"Keywords", Value:"博客,个人博客,梦想星辰大海"})
				siteInfoItem = append(siteInfoItem, &m_configItem.ConfigItem{GroupID:siteInfo.ID, Name:"站点描述", Key:"Description", Value:"个人的博客，记录所思所想。"})
				siteInfoItem = append(siteInfoItem, &m_configItem.ConfigItem{GroupID:siteInfo.ID, Name:"底部信息", Key:"Footer", Value:"© "+time.Now().Format("2006")})
				siteInfoItem = append(siteInfoItem, &m_configItem.ConfigItem{GroupID:siteInfo.ID, Name:"统计代码", Key:"Statistical", Value:""})
				if _, err := memoryDao.Insert(siteInfoItem); err != nil {
					a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
					os.Exit(1)
				}
			}

			//导出数据
			var tables []*schemas.Table
			tables, err = dao.Dao.DBMetas()
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			var f *os.File
			f, err = os.OpenFile("database/init.sql", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			defer func() {
				_ = f.Close()
			}()
			err = s_services.DumpDB(memoryDao, tables, f, s_services.DUMP_DB_DATA)
			if err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			if err := memoryDao.Close(); err != nil {
				a_boot.Logger.ErrorF("导出数据库到./database/init.sql失败: %s", err)
				os.Exit(1)
			}
			a_boot.Logger.InfoF("导出数据库到 database/init.sql 成功")
		},
	}
	dbCmd.AddCommand(dumpInitCmd)
}

