package models

import "time"

//id
type IDField struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
}

func (this IDField) GetID() int {
	return this.ID
}

//创建时间
type CreatedAtField struct {
	CreatedAt time.Time `xorm:"DATETIME created"`
}

func (this CreatedAtField) CreatedAtText() string {
	return this.CreatedAt.Format("2006-01-02 15:04:05")
}

//更新时间
type UpdatedAtField struct {
	UpdatedAt time.Time `xorm:"DATETIME updated"`
}

func (this UpdatedAtField) UpdatedAtText() string {
	return this.UpdatedAt.Format("2006-01-02 15:04:05")
}

//软删除时间
type DeletedAtField struct {
	DeletedAt time.Time `xorm:"DATETIME deleted"`
}

func (this DeletedAtField) DeletedAtText() string {
	return this.DeletedAt.Format("2006-01-02 15:04:05")
}

//公共字段
type Field struct {
	IDField        `xorm:"extends"`
	CreatedAtField `xorm:"extends"`
	UpdatedAtField `xorm:"extends"`
	DeletedAtField `xorm:"extends"`
}
