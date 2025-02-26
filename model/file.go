package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName      string `gorm:"type:string;not null;"`
	FilePath      string `gorm:"type:string;unique;not null;"`
	IsChunk       bool   `gorm:"type:boolean;not null;default:1"` // 是否分块
	MD5           string `gorm:"type:string;"`
	Size          int64  `gorm:"type:bigint;"`
	Ext           string `gorm:"type:string;"`
	Description   string `gorm:"type:string;"`
	DownloadTimes int64  `gorm:"type:bigint;"`
	AuthorID      uint   `gorm:"index"`
	Author        User   `gorm:"foreignKey:AuthorID;References:ID"`
}

type Chunk struct {
	gorm.Model
	FileName   string `gorm:"type:string;not null;"`
	FilePath   string `gorm:"type:string;unique;not null;"`
	ChunkIndex int    `gorm:"type:int;not null;index"`
	MD5        string `gorm:"type:string;index"`
	Size       int64  `gorm:"type:bigint;"`
	Ext        string `gorm:"type:string;"`
	AuthorID   uint   `gorm:"index"`
	Author     User   `gorm:"foreignKey:AuthorID;References:ID"`
}
