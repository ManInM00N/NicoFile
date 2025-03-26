package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	Name string `gorm:"type:string;not null;"` // 图片名称
	Path string `gorm:"type:string;not null;"` // 图片路径

	AuthorID uint `gorm:"index"`
	Author   User `gorm:"foreignKey:AuthorID;References:ID"`
}
