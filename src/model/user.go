package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	Id          uint32     `gorm:"primary_key" json:"id"`
	AppId       uint32     `gorm:"index" json:"app_id"`
	Nickname    string     `gorm:"size:30;index" json:"nickname"`
	Avatar      string     `json:"avatar"`
	Gender      uint8      `json:"gender"` // refer to constant/typ.go Gender*
	Birthday    *time.Time `grom:"type:date" json:"birthday"`
	Country     string     `gorm:"size:16" json:"country"`
	City        string     `gorm:"size:100" json:"city"`
	Address     string     `gorm:"size:255" json:"address"`
	CheckinDays uint32     `json:"checkin_days"` // 连续签到天数
	LastCheckin *time.Time `json:"last_checkin"` // 最后签到时间
	Gold        uint32     `json:"gold"`         // 金币余额
	InGold      uint32     `json:"in_gold"`      // 获得的金币
	OutGold     uint32     `json:"out_gold"`     // 已花费的金币
	Status      uint32     `json:"status"`
	//Auths       []*UserAuth      `json:"auths"`
	//Vips        []*UserVip       `json:"vips"`
	//Privileges  []*UserPrivilege `json:"privileges"`
	Longitude float32   `json:"longitude"` // 评论发布的经度
	Latitude  float32   `json:"latitude"`  // 评论发布的纬度
	CreatedIp string    `gorm:"size:96" json:"created_ip"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	UserStatusLoginBanned   = (1 << 0) // 用户状态：禁止登录
	UserStatusCommentBanned = (1 << 1) // 用户状态：禁止发言
	UserStatusWhiteList     = (1 << 2) // 用户状态：白名单

	VipStatusVip     = 2 // vip状态
	VipStatusExpired = 1 // 过期的vip
	VipStatusUnknown = 0 // 不是vip
)

func (User) TableName() string {
	return "user"
}

func initUser(db *gorm.DB) error {
	var err error

	if db.HasTable(&User{}) {
		err = db.AutoMigrate(&User{}).Error
	} else {
		err = db.CreateTable(&User{}).Error
	}
	return err
}

func dropUser(db *gorm.DB) {
	db.DropTableIfExists(&User{})
}
