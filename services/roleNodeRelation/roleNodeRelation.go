package s_roleNodeRelation

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/node"
	"github.com/buexplain/go-blog/models/roleNodeRelation"
)

type Relation struct {
	m_node.Node `xorm:"extends"`
	Checked     bool `xorm:"-"`
}

type RelationList []*Relation

func (this RelationList) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

//获取角色节点
func GetRelation(roleID int) (RelationList, error) {
	//获取所有的节点
	allNode := make(RelationList, 0)
	err := dao.Dao.Table("Node").Desc("SortID").Find(&allNode)
	if err != nil {
		return nil, err
	}

	//获取角色拥有的节点
	roleNode := make(m_roleNodeRelation.List, 0)
	err = dao.Dao.Table("RoleNodeRelation").Where("RoleID=?", roleID).Find(&roleNode)
	if err != nil {
		return nil, err
	}
	if len(roleNode) > 0 {
		for _, node := range allNode {
			node.Checked = roleNode.HasNodeID(node.ID)
		}
	}

	return allNode, nil
}

//设置角色的节点
func SetRelation(roleID int, nodeID []int) error {
	//开启事务
	session := dao.Dao.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	//先删除已有关系
	r := &m_roleNodeRelation.RoleNodeRelation{}
	r.RoleID = roleID
	_, err := session.Delete(r)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return err
		}
		return err
	}

	//插入新关系
	result := make(m_roleNodeRelation.List, 0, len(nodeID))
	for _, v := range nodeID {
		result = append(result, m_roleNodeRelation.RoleNodeRelation{RoleID: roleID, NodeID: v})
	}
	_, err = session.Insert(&result)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return err
		}
		return err
	} else {
		if err := session.Commit(); err != nil {
			return err
		}
		return nil
	}
}

//节点角色id
type NodeRoleID int

type NodeRoleIDList []NodeRoleID

//获取节点的角色id列表
func GetRoleIDByNodeURL(nodeURL string, method ...string) (NodeRoleIDList, error) {
	result := make(NodeRoleIDList, 0)
	session := dao.Dao.Table("`Node`").Where("`Node`.`URL`=?", nodeURL)
	session.Join("INNER", "`RoleNodeRelation`", "`RoleNodeRelation`.NodeID = `Node`.ID")
	session.Select("`RoleNodeRelation`.`RoleID` as `NodeRoleID`")
	if len(method) > 0 {
		session.Where("`Node`.`Methods` LIKE ?", "%"+method[0]+"%")
	}
	err := session.Find(&result)
	return result, err
}
