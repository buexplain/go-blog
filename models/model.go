package models

import "time"

//公共字段
type Field struct {
	ID        int       `xorm:"not null pk autoincr INTEGER"`
	CreatedAt time.Time `xorm:"DATETIME created"`
	UpdatedAt time.Time `xorm:"DATETIME updated"`
	DeletedAt time.Time `xorm:"DATETIME deleted"`
}

func (this Field) GetID() int {
	return this.ID
}

func (this Field) CreatedAtText() string {
	return this.CreatedAt.Format("2006-01-02 15:04:05")
}

func (this Field) UpdatedAtText() string {
	return this.UpdatedAt.Format("2006-01-02 15:04:05")
}

func (this Field) DeletedAtText() string {
	return this.DeletedAt.Format("2006-01-02 15:04:05")
}
