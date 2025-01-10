package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:string;unique;not null;"`
	Password string
}
