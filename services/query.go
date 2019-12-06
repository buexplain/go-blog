package s_services

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/util"
	"github.com/buexplain/go-fool"
	"xorm.io/xorm"
	"time"
	"xorm.io/core"
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
		this.tableInfo, this.Error = m_util.GetTableInfo(dao.Dao, this.tableName)
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

//查询并且统计行数
func (this *Query) FindAndCount(rowsSlicePtr interface{}, counter *int64, condiBean ...interface{}) {
	if this.Error == nil {
		c, err := this.Finder.FindAndCount(rowsSlicePtr, condiBean...)
		if err != nil {
			this.Error = err
		}else {
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
