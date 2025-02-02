package model

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"type:string;unique;not null;" json:"username"`
	Password   string `json:"-"`
}
