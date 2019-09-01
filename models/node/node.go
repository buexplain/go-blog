package m_node

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
	"strings"
)

//rbac 节点表
type Node struct {
	models.Field `xorm:"extends"`
	//父id
	Pid          int    `xorm:"INTEGER"`
	//节点名称
	Name         string `xorm:"TEXT"`
	//节点路径
	URL          string `xorm:"TEXT"`
	//节点路径允许的方法
	Methods       string `xorm:"TEXT"`
	//是否为后台导航菜单
	IsMenu       int    `xorm:"INTEGER"`
	//排序id
	SortID       int    `xorm:"INTEGER"`
}

func (this Node) HasMethod(method string) bool {
	if this.Methods == method {
		return true
	}
	if strings.Index(this.Methods, method) == -1 {
		return false
	}
	return true
}

func (this *Node) Parse() {
	if this.IsMenu == 0 {
		this.IsMenu = IsMenuNo
	}else {
		this.IsMenu = IsMenuYes
	}
	this.URL = strings.Trim(this.URL, " ")
}

func (this *Node) Insert() (affected int64, err error) {
	this.Parse()
	return dao.Dao.Insert(this)
}

func (this *Node) Update() (affected int64, err error) {
	this.Parse()
	return dao.Dao.ID(this.ID).MustCols("Pid", "Methods").Update(this)
}

type List []Node

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func GetALL() (List, error) {
	result := make(List, 0)
	err := dao.Dao.Table("Node").Desc("SortID").Find(&result)
	return result, err
}

func GetMenu() (List, error) {
	result := make(List, 0)
	err := dao.Dao.Table("Node").Where("IsMenu=?", IsMenuYes).Desc("SortID").Find(&result)
	return result, err
}

func GetMenuByUserID(userID int) (List, error) {
	result := make(List, 0)
	mod := dao.Dao.Table("Node").
		Join("INNER", "`RoleNodeRelation`", "`Node`.`ID` = `RoleNodeRelation`.`NodeID`").
		Join("INNER", "`UserRoleRelation`", "`UserRoleRelation`.`RoleID` = `RoleNodeRelation`.`RoleID`").
		Where("UserRoleRelation.UserID=?", userID).Where("Node.IsMenu=?", IsMenuYes)
	err := mod.Find(&result)
	return result, err
}