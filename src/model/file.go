package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type File struct {
	Id          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	ContentType uint32     `json:"content_type"`
	RelPath     string     `gorm:"size:255" json:"rel_path"`
	RelHash     uint32     `gorm:"index" json:"rel_hash"`
	SourceUrl   string     `gorm:"size:512" json:"source_url"`
	SourceHash  uint32     `gorm:"index" json:"source_hash"`
	Ext         string     `gorm:"size:16" json:"ext"`
	Filename    string     `gorm:"size:128" json:"filename"`
	Filesize    int64      `json:"filesize"`
	LastAccess  *time.Time `json:"last_access"`
	AccessCount uint32     `json:"access_count"`
	Tag         string     `gorm:"size:64" json:"tag"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (File) TableName() string {
	return "file"
}

func initFile(db *gorm.DB) error {
	var err error

	if db.HasTable(&File{}) {
		err = db.AutoMigrate(&File{}).Error
	} else {
		err = db.CreateTable(&File{}).Error
	}

	return err
}

func dropFile(db *gorm.DB) {
	db.DropTableIfExists(&File{})
}
