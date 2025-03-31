package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ArticleID uint `gorm:"not null;index"`
	AuthorId  uint
	Author    User      // 评论者名称
	Content   string    // 评论内容
	IP        string    // 评论者IP(用于反垃圾)
	Status    string    // 评论状态: approved, pending, spam
	ParentID  *uint     // 用于实现回复功能，指向父评论ID
	Article   Article   `gorm:"references:ID"`
	Replies   []Comment `gorm:"foreignkey:ParentID"` // 自引用关系，用于回复
}
