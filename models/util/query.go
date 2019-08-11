package m_util

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-fool"
	"github.com/go-xorm/xorm"
	"strconv"
	"time"
	"xorm.io/core"
)

//列表查询
type Query struct {
	//请求上下文
	ctx               *fool.Ctx
	//表名称
	tableName         string
	//关键词字段的key
	KeywordFieldKey   string
	//关键词的值的key
	KeywordValueKey   string
	//时间字段的key
	TimeFieldKey      string
	//开始时间的值的key
	TimeStartValueKey string
	//结束时间的值的key
	TimeEndValueKey   string
	//表信息
	tableInfo         *core.Table
	//统计
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
		size = append(size, 20)
	}
	page := this.ctx.Request().QueryInt("page", 1)
	this.ctx.Response().Assign("limit", size)
	this.ctx.Response().Assign("page", page)
	offset := (page - 1) * size[0]
	this.Finder.Limit(size[0], offset)
	return this
}

//关键词like
func (this *Query) WhereKeyword() *Query {
	keywordField := this.ctx.Request().Query(this.KeywordFieldKey)
	keywordValue := this.ctx.Request().Query(this.KeywordValueKey)
	this.ctx.Response().Assign(this.KeywordFieldKey, keywordField)
	this.ctx.Response().Assign(this.KeywordValueKey, keywordValue)
	if keywordField == "" || keywordValue == "" {
		return this
	}
	tableInfo := this.TableInfo()
	if tableInfo == nil {
		return this
	}
	column := tableInfo.GetColumn(keywordField)
	if column == nil {
		return this
	}
	if column.IsAutoIncrement {
		this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), keywordValue)
	} else {
		if column.SQLType.IsBlob() {
			b, _ := strconv.ParseBool(keywordValue)
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), b)
		} else if column.SQLType.IsNumeric() {
			this.Finder.Where(fmt.Sprintf("`%s`.`%s`=?", this.tableName, column.Name), keywordValue)
		} else {
			this.Finder.Where(fmt.Sprintf("`%s`.`%s` LIKE ?", this.tableName, column.Name), fmt.Sprintf("%s%s%s", "%", keywordValue, "%"))
		}
	}
	return this
}

//时间范围
func (this *Query) WhereTime() *Query {
	timeField := this.ctx.Request().Query(this.TimeFieldKey)
	timeStartValue := this.ctx.Request().Query(this.TimeStartValueKey)
	timeEndValue := this.ctx.Request().Query(this.TimeEndValueKey)
	this.ctx.Response().Assign(this.TimeFieldKey, timeField)
	this.ctx.Response().Assign(this.TimeStartValueKey, timeStartValue)
	this.ctx.Response().Assign(this.TimeEndValueKey, timeEndValue)
	if timeField == "" || (timeStartValue == "" && timeEndValue == "") {
		return this
	}
	tableInfo := this.TableInfo()
	if tableInfo == nil {
		return this
	}
	column := tableInfo.GetColumn(timeField)
	if column == nil {
		return this
	}
	if timeStartValue != "" {
		this.Finder.Where(fmt.Sprintf("`%s`.`%s`>?", this.tableName, column.Name), timeStartValue)
	}
	if timeEndValue != "" {
		this.Finder.Where(fmt.Sprintf("`%s`.`%s`<?", this.tableName, column.Name), timeEndValue)
	}
	return this
}


func NewQuery(tableName string, ctx *fool.Ctx) *Query {
	tmp := new(Query)
	tmp.ctx = ctx
	tmp.tableName = tableName
	tmp.KeywordFieldKey = "keywordField"
	tmp.KeywordValueKey = "keywordValue"
	tmp.TimeFieldKey = "timeField"
	tmp.TimeStartValueKey = "timeStartValue"
	tmp.TimeEndValueKey = "timeEndValue"
	tmp.Finder = dao.Dao.Table(tmp.tableName)
	return tmp
}
