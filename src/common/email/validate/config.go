package validate

import "time"

// Config 邮箱验证的配置参数
type Config struct {
	Expire  time.Duration // 过期的持续时间
	CodeLen int           // 验证码的长度
}

// DataItem 存储验证信息的数据项
type DataItem struct {
	ID         int64         // 唯一标识
	Email      string        // 邮箱
	Code       string        // 验证码
	CreateTime time.Time     // 存储时间
	Expire     time.Duration // 过期的持续时间
}