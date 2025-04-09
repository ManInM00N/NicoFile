package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title    string    `gorm:"type:string;not null;" json:"title"` // 文章标题
	Content  string    `gorm:"type:text;not null" json:"content"`
	AuthorID uint      `gorm:"index" json:"author_id"` // 文章作者ID
	Author   User      `gorm:"foreignKey:AuthorID;References:ID"`
	View     int64     `gorm:"default:0" json:"view"`
	Like     int64     `gorm:"default:0" json:"like"`
	Cover    string    `gorm:"type:string;" json:"cover"`
	Comments []Comment `gorm:"foreignKey:ArticleID;References:ID" json:"comments"` // 一对多关系: 一篇文章有多条评论
}
