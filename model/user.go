package model

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"type:string;unique;not null;index" json:"username"`
	Password   string `json:"-"`
	Priority   int    `gorm:"type:int;not null;default:0" json:"-"`
}
