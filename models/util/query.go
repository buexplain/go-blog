package m_util

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-fool"
	"github.com/go-xorm/xorm"
	"time"
	"xorm.io/core"
)

type where []string

func (this where) Get(index int) string {
	if index+1 <= len(this) {
		return this[index]
	}
	return ""
}

//列表查询
type Query struct {
	//请求上下文
	ctx               *fool.Ctx
	//表名称
	tableName         string
	//表信息
	tableInfo         *core.Table
	//统计器
	Counter           *xorm.Session
	//查询器
	Finder            *xorm.Session
	//错误值
	Error             error
}

//返回表信息
func (this *Query) TableInfo() *core.Table {
	if this.Error == nil && this.tableInfo == nil {
		this.tableInfo, this.Error = GetTableInfo(dao.Dao, this.tableName)
		if this.Error == nil {
			this.tableName = this.tableInfo.Name
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
		this.Error = this.Finder.Find(rowsSlicePtr)
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
			} else if tableInfo.GetColumn("ID") != nil {
				type Tmp struct {
					ID int `xorm:"INTEGER"`
				}
				total, this.Error = this.Counter.Count(new(Tmp))
			} else {
				//如果报这个错误，请在此处添加一个比较通用的字段
				this.Error = fmt.Errorf("not hit table field: %s", this.tableName)
			}
		}
	}
	return total
}

//分页大小
func (this *Query) Limit(size ...int) *Query {
	if len(size) == 0 {
		size = append(size, 10)
	}
	page := this.ctx.Request().QueryInt("page", 1)
	this.ctx.Response().Assign("limit", size)
	this.ctx.Response().Assign("page", page)
	offset := (page - 1) * size[0]
	this.Finder.Limit(size[0], offset)
	return this
}

//筛选，这些条件都是 and 的
func (this *Query) Screen() *Query {
	//等于查询的参数
	eqFieldSlice := where(this.ctx.Request().QuerySlice("eqField[]"))
	eqValueSlice := where(this.ctx.Request().QuerySlice("eqValue[]"))
	this.ctx.Response().Assign("eqField", eqFieldSlice)
	this.ctx.Response().Assign("eqValue", eqValueSlice)

	//大于等于查询的参数
	geFieldSlice := where(this.ctx.Request().QuerySlice("geField[]"))
	geValueSlice := where(this.ctx.Request().QuerySlice("geValue[]"))
	this.ctx.Response().Assign("geField", geFieldSlice)
	this.ctx.Response().Assign("geValue", geValueSlice)

	//小于等于查询的参数
	leFieldSlice := where(this.ctx.Request().QuerySlice("leField[]"))
	leValueSlice := where(this.ctx.Request().QuerySlice("leValue[]"))
	this.ctx.Response().Assign("leField", leFieldSlice)
	this.ctx.Response().Assign("leValue", leValueSlice)

	//like查询
	likeFieldSlice := where(this.ctx.Request().QuerySlice("likeField[]"))
	likeValueSlice := where(this.ctx.Request().QuerySlice("likeValue[]"))
	this.ctx.Response().Assign("likeField", likeFieldSlice)
	this.ctx.Response().Assign("likeValue", likeValueSlice)

	tableInfo := this.TableInfo()
	if tableInfo == nil {
		return this
	}

	//等于查询
	if len(eqFieldSlice) > 0 && (len(eqFieldSlice) == len(eqValueSlice)) {
		for index, field := range eqFieldSlice {
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			value := eqValueSlice[index]
			if value == "" {
				continue
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), value)
		}
	}

	//大于等于查询
	if len(geFieldSlice) > 0 && (len(geFieldSlice) == len(geValueSlice)) {
		for index, field := range geFieldSlice {
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			value := geValueSlice[index]
			if value == "" {
				continue
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`>=?", this.tableName, column.Name), value)
		}
	}

	//小于等于查询
	if len(leFieldSlice) > 0 && (len(leFieldSlice) == len(leValueSlice)) {
		for index, field := range leFieldSlice {
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			value := leValueSlice[index]
			if value == "" {
				continue
			}
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`<=?", this.tableName, column.Name), value)
		}
	}

	//like查询
	if len(likeFieldSlice) > 0 && (len(likeFieldSlice) == len(likeValueSlice)) {
		for index, field := range likeFieldSlice {
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			value := likeValueSlice[index]
			if value == "" {
				continue
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
