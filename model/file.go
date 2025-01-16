package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName    string `gorm:"type:string;not null;"`
	FilePath    string `gorm:"type:string;unique;not null;"`
	IsChunk     bool   `gorm:"type:boolean;not null;"` // 是否分块
	MD5         string `gorm:"type:string;"`
	Size        int64  `gorm:"type:bigint;"`
	Ext         string `gorm:"type:string;"`
	Description string `gorm:"type:string;"`
}
