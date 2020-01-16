package s_services

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/util"
	"github.com/buexplain/go-fool"
	"github.com/buexplain/go-fool/errors"
	"strings"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

type where []string

func (this where) Get(index int) string {
	if index+1 <= len(this) && index >= 0 {
		return this[index]
	}
	return ""
}

//查询构造器
type Query struct {
	//请求上下文
	ctx *fool.Ctx
	//表名称
	tableName string
	//表信息
	tableInfo *core.Table
	//统计器
	Counter *xorm.Session
	//查询器
	Finder *xorm.Session
	//错误值
	Error error
}

//返回表信息
func (this *Query) TableInfo() *core.Table {
	if this.Error == nil && this.tableInfo == nil {
		this.tableInfo, this.Error = m_util.GetTableInfo(dao.Dao, this.tableName)
		if this.Error == nil {
			this.tableName = this.tableInfo.Name
		}else {
			this.Error = errors.MarkServer(this.Error)
		}
	}
	return this.tableInfo
}

//查询
func (this *Query) Find(rowsSlicePtr interface{}) {
	if this.Counter == nil {
		this.Counter = this.Finder.Clone()
	}
	if this.Error == nil {
		this.Error = errors.MarkServer(this.Finder.Find(rowsSlicePtr))
	}
}

//分页统计
func (this *Query) Count() int64 {
	if this.Counter == nil {
		this.Counter = this.Finder.Clone()
	}
	var total int64
	if this.Error == nil {
		tableInfo := this.TableInfo()
		if tableInfo != nil {
			if tableInfo.GetColumn("DeletedAt") != nil {
				type Tmp struct {
					DeletedAt time.Time `xorm:"DATETIME deleted"`
				}
				total, this.Error = this.Counter.Count(new(Tmp))
				this.Error = errors.MarkServer(this.Error)
			} else if tableInfo.GetColumn("ID") != nil {
				type Tmp struct {
					ID int `xorm:"INTEGER"`
				}
				total, this.Error = this.Counter.Count(new(Tmp))
				this.Error = errors.MarkServer(this.Error)
			} else {
				//如果报这个错误，请在此处添加一个比较通用的字段
				this.Error = errors.MarkServer(fmt.Errorf("not hit table field: %s", this.tableName))
			}
		}
	}
	return total
}

//查询并且统计行数
func (this *Query) FindAndCount(rowsSlicePtr interface{}, counter *int64, condiBean ...interface{}) {
	if this.Error == nil {
		c, err := this.Finder.FindAndCount(rowsSlicePtr, condiBean...)
		if err != nil {
			this.Error = errors.MarkServer(err)
		} else {
			*counter = c
		}
	}
}

//分页大小
func (this *Query) Limit() *Query {
	page := this.ctx.Request().QueryInt("page", 1)
	limit := this.ctx.Request().QueryInt("limit", 10)
	this.ctx.Response().Assign("limit", limit)
	this.ctx.Response().Assign("page", page)
	offset := (page - 1) * limit
	this.Finder.Limit(limit, offset)
	return this
}

func (this *Query) Where() *Query {
	tableInfo := this.TableInfo()
	if tableInfo == nil {
		return this
	}

	//等于查询
	whereEq := this.ctx.Request().QueryMap("whereEq")
	if whereEq != nil {
		for field, value := range whereEq {
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			if value == "" {
				continue
			}
			if !this.ctx.Request().IsAjax() {
				this.ctx.Response().Assign(field, value)
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), value)
		}
	}
	//大于等于查询的参数
	whereGe := this.ctx.Request().QueryMap("whereGe")
	if whereGe != nil {
		for field, value := range whereGe {
			if value == "" {
				continue
			}
			//判断字段是否存在
			column := tableInfo.GetColumn(strings.TrimLeft(field, "Ge"))
			if column == nil {
				continue
			}
			if !this.ctx.Request().IsAjax() {
				this.ctx.Response().Assign(field, value)
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`>=?", this.tableName, column.Name), value)
		}
	}
	//小于等于查询
	whereLe := this.ctx.Request().QueryMap("whereLe")
	if whereLe != nil {
		for field, value := range whereLe {
			if value == "" {
				continue
			}
			//判断字段是否存在
			column := tableInfo.GetColumn(strings.TrimLeft(field, "Le"))
			if column == nil {
				continue
			}
			if !this.ctx.Request().IsAjax() {
				this.ctx.Response().Assign(field, value)
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`<=?", this.tableName, column.Name), value)
		}
	}
	//模糊查询
	whereLike := this.ctx.Request().QueryMap("whereLike")
	if whereLike != nil {
		for field, value := range whereLike {
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			if value == "" {
				continue
			}
			if !this.ctx.Request().IsAjax() {
				this.ctx.Response().Assign(field, value)
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s` LIKE ?", this.tableName, column.Name), fmt.Sprintf("%s%s%s", "%", value, "%"))
		}
	}
	return this
}

func NewQuery(tableName string, ctx *fool.Ctx) *Query {
	tmp := new(Query)
	tmp.ctx = ctx
	tmp.tableName = tableName
	tmp.Finder = dao.Dao.Table(tmp.tableName)
	return tmp
}
