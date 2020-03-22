package m_roleNodeRelation

//rbac 角色与节点的关系表
type RoleNodeRelation struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	//角色id
	RoleID int `xorm:"INTEGER"`
	//节点
	NodeID int `xorm:"INTEGER"`
}

type List []RoleNodeRelation

func (this List) HasNodeID(nodeID int) bool {
	for _, v := range this {
		if v.NodeID == nodeID {
			return true
		}
	}
	return false
}
