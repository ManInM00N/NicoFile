package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title    string `gorm:"type:string;not null;"` // 文章标题
	Content  string `gorm:"type:text;not null"`
	AuthorID uint   `gorm:"index"`
	Author   User   `gorm:"foreignKey:AuthorID;References:ID"`
	View     int64  `gorm:"default:0"`
}
