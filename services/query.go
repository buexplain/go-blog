package s_services

import (
	"fmt"
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/helpers"
	m_models "github.com/buexplain/go-blog/models"
	"github.com/buexplain/go-fool"
	"strings"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
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
	tableInfo *schemas.Table
	//是否注入where条件到模板
	assign bool
	//是否调用了where方法
	condition int8
	//查询器
	Finder *xorm.Session
	//错误值
	Error error
}

//返回表信息
func (this *Query) TableInfo() *schemas.Table {
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
	if this.condition == 2 {
		h_boot.Logger.Warning("where condition are gone")
	}
	//调用一次当前方法，如果条件就绪，则标记为2，表示条件已经消费
	if this.condition == 1 {
		this.condition = 2
	}
	if this.Error == nil {
		this.Error = this.Finder.Find(rowsSlicePtr)
	}
}

//分页统计
func (this *Query) Count() int64 {
	if this.condition == 2 {
		h_boot.Logger.Warning("where condition are gone")
	}
	//调用一次当前方法，如果条件就绪，则标记为2，表示条件已经消费
	if this.condition == 1 {
		this.condition = 2
	}
	var total int64
	if this.Error == nil {
		tableInfo := this.TableInfo()
		if tableInfo != nil {
			if tableInfo.GetColumn("CreatedAt") != nil {
				type Tmp struct {
					CreatedAt m_models.Time `xorm:"DATETIME created"`
				}
				total, this.Error = this.Finder.Count(new(Tmp))
			} else if tableInfo.GetColumn("ID") != nil {
				type Tmp struct {
					ID int `xorm:"INTEGER"`
				}
				total, this.Error = this.Finder.Count(new(Tmp))
			} else {
				//如果报这个错误，请在此处添加一个比较通用的字段
				this.Error = code.NewF(code.SERVER, "not hit table field: %s", this.tableName)
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
		} else {
			*counter = c
		}
	}
}

//分页大小
func (this *Query) Limit() *Query {
	page := this.ctx.Request().QueryInt("page", 1)
	limit := this.ctx.Request().QueryInt("limit", 10)
	if this.assign {
		this.ctx.Response().Assign("limit", limit)
		this.ctx.Response().Assign("page", page)
	}
	offset := (page - 1) * limit
	this.Finder.Limit(limit, offset)
	return this
}

func (this *Query) Where() *Query {
	//调用一次where，则条件设置为1，表示条件就绪
	this.condition = 1
	tableInfo := this.TableInfo()
	if tableInfo == nil {
		return this
	}

	//等于查询
	whereEq := this.ctx.Request().QueryMap("whereEq")
	if whereEq != nil {
		for field, value := range whereEq {
			if value == "" {
				continue
			}
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			if this.assign {
				this.ctx.Response().Assign(field, value)
			}
			if column.SQLType.IsTime() {
				//如果是时间字段，则将时间解析为UTC时间
				if t, err := helpers.ParseInLocation(value, time.Local); err != nil {
					h_boot.Logger.WarningF("parse where condition %s error: %s", column.Name, err)
				} else {
					this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), t.UTC().Format("2006-01-02 15:04:05"))
				}
			} else {
				this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), value)
			}
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
			if this.assign {
				this.ctx.Response().Assign(field, value)
			}
			if column.SQLType.IsTime() {
				//如果是时间字段，则将时间解析为UTC时间
				if t, err := helpers.ParseInLocation(value, time.Local); err != nil {
					h_boot.Logger.WarningF("parse where condition %s error: %s", column.Name, err)
				} else {
					this.Finder.Where(fmt.Sprintf("`%s`.`%s`>=?", this.tableName, column.Name), t.UTC().Format("2006-01-02 15:04:05"))
				}
			} else {
				this.Finder.Where(fmt.Sprintf("`%s`.`%s`>=?", this.tableName, column.Name), value)
			}
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
			if this.assign {
				this.ctx.Response().Assign(field, value)
			}
			if column.SQLType.IsTime() {
				//如果是时间字段，则将时间解析为UTC时间
				if t, err := helpers.ParseInLocation(value, time.Local); err != nil {
					h_boot.Logger.WarningF("parse where condition %s error: %s", column.Name, err)
				} else {
					this.Finder.Where(fmt.Sprintf("`%s`.`%s`<=?", this.tableName, column.Name), t.UTC().Format("2006-01-02 15:04:05"))
				}
			} else {
				this.Finder.Where(fmt.Sprintf("`%s`.`%s`<=?", this.tableName, column.Name), value)
			}
		}
	}
	//模糊查询
	whereLike := this.ctx.Request().QueryMap("whereLike")
	if whereLike != nil {
		for field, value := range whereLike {
			if value == "" {
				continue
			}
			//判断字段是否存在
			column := tableInfo.GetColumn(field)
			if column == nil {
				continue
			}
			if this.assign {
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
	//是ajax请求，或者要求按json进行响应的路由，则不注入参数到模板
	tmp.assign = !(ctx.Request().IsAjax() || (ctx.Route() != nil && ctx.Route().HasLabel("json")))
	tmp.condition = 0
	return tmp
}
