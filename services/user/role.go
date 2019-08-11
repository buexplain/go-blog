package s_user

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	m_role "github.com/buexplain/go-blog/models/role"
	m_userRoleRelation "github.com/buexplain/go-blog/models/userRoleRelation"
)

type Tree struct {
	Title string  `json:"title"`
	ID int  `json:"id"`
	Pid int  `json:"pid"`
	Children []*Tree  `json:"children"`
	Href string  `json:"href"`
	Spread bool  `json:"spread"`
	Checked bool  `json:"checked"`
	Disabled bool  `json:"disabled"`
}

type TreeList []*Tree

func (this TreeList) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func GetRoleTree(userID int) (TreeList, error) {
	//获取所有的角色
	allRole := make([]m_role.Role, 0)
	err := dao.Dao.Table("Role").Desc("SortID").Find(&allRole)
	if err != nil {
		return nil, err
	}

	//获取用户拥有的角色
	userRole := make(m_userRoleRelation.List, 0)
	err = dao.Dao.Table("UserRoleRelation").Find(&userRole)
	if err != nil {
		return nil, err
	}

	//转化成一棵树
	treeMap := make(map[int]*Tree)
	treeList := make([]*Tree, 0, len(allRole))
	for _, role := range allRole {
		t := &Tree{
			Title:role.Name,
			ID:role.ID,
			Pid:role.Pid,
			Children:make(TreeList, 0),
			Href:"",
			Spread:true,
			Checked:userRole.HasRoleID(role.ID),
			Disabled:false,
		}
		treeMap[role.ID] = t
		treeList = append(treeList, t)
	}

	result := make(TreeList, 0, len(treeList))
	for _, v := range treeList {
		if v.Pid == 0 {
			result = append(result, v)
		}else {
			treeMap[v.Pid].Children = append(treeMap[v.Pid].Children, v)
		}
	}
	return result, nil
}