package constant

import "time"

const (
	// 验证码配置
	DefaultExpire     = time.Hour * 2    // 默认过期时间（2个小时）
	DefaultCodeLen    = 6                // 默认验证码的长度
	DefaultGCInterval = time.Second * 60 // 默认验证信息的GC间隔
)

const (
	FileContentTypePicture = 1 //文件类型，图片
)

const (
	OsTypeUnknown = 0 // 未知类型
	OsTypeAndroid = 1 // android
	OsTypeIos     = 2 // ios
)