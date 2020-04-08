package e_syncRbacNode

import (
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/dao"
	m_node "github.com/buexplain/go-blog/models/node"
	m_role "github.com/buexplain/go-blog/models/role"
	s_roleNodeRelation "github.com/buexplain/go-blog/services/roleNodeRelation"
	"github.com/buexplain/go-event"
)

//事件名称
const EVENT_NAME = "syncRbacNode"

//事件监听者
type EventListener struct {
}

func (this *EventListener) Handle(e *event.Event) {
	superRoleID, ok := e.Data.(int)
	if !ok || superRoleID == 0 {
		h_boot.Logger.ErrorF("同步超级角色节点失败，超级角色ID错误：%+v", e.Data)
		return
	}
	result := new(m_role.Role)
	result.ID = superRoleID
	if has, err := dao.Dao.Get(result); err != nil {
		h_boot.Logger.ErrorF("同步超级角色节点失败: %s", err)
		return
	} else if !has {
		h_boot.Logger.InfoF("同步超级角色节点失败，没有找到角色: %d", result.ID)
		return
	}
	if l, err := m_node.GetALL(); err != nil {
		h_boot.Logger.ErrorF("同步超级角色节点失败: %s", err)
		return
	} else {
		nodeID := make([]int, 0, len(l))
		for _, v := range l {
			nodeID = append(nodeID, v.ID)
		}
		if len(nodeID) > 0 {
			if err := s_roleNodeRelation.SetRelation(result.ID, nodeID); err != nil {
				h_boot.Logger.ErrorF("同步超级角色节点失败: %s", err)
			} else {
				h_boot.Logger.InfoF("同步超级角色节点成功: %d", result.ID)
			}
		} else {
			h_boot.Logger.Info("同步超级角色节点失败: 没有找到任何权限节点")
		}
	}
}
