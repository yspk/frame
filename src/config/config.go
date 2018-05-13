package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	BindAddr           string    `json:"bind_addr"`
	ProductionEnv      bool      `json:"production_env"`
	Root               string    `json:"ims_root"`
	LogRoot            string    `json:"log_root"`
	StorageRoot        string    `json:"storage_root"`
	DBName             string    `json:"db_name"`
	DBSource           string    `json:"db_source"`
	LoggerLevel        uint8     `json:"logger_level"`
	OrmLogEnabled      bool      `json:"orm_log_enabled"`
	CacheRedisAddr     string    `json:"cache_redis_addr"`
	CacheRedisPassword string    `json:"cache_redis_password"`
	ValidateEmail      ValEmail  `json:"validate_email"`
	ChuangLanConfig    ChuangLan `json:"chuang_lan_config"`
}

type ValEmail struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type ChuangLan struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Template string `json:"template"`
	Uid      string `json:"uid"`
}

var c config

func init() {
	c.BindAddr = "http://192.168.3.165:8555"
	c.ProductionEnv = false
	c.Root = "../"
	c.LogRoot = "../log/"
	c.StorageRoot = "../data/"
	c.DBName = "mysql"
	c.DBSource = "root:toor@tcp(localhost:3306)/front?charset=utf8&parseTime=True&loc=Local"
	c.LoggerLevel = 1
	c.OrmLogEnabled = true
	c.CacheRedisAddr = ""
	c.CacheRedisPassword = ""
	c.ValidateEmail.Account = "ys@baoquan.com"
	c.ValidateEmail.Name = "yusheng"
	c.ValidateEmail.Password = "QAZxsw2"
	c.ValidateEmail.Host = "smtp.exmail.qq.com:25"
}

func LoadConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var ctmp config
	err = json.Unmarshal(b, &ctmp)
	if err != nil {
		return err
	}

	c = ctmp
	return nil
}

func IsProductionEnv() bool {
	return c.ProductionEnv
}

func GetImsRoot() string {
	return c.Root
}

func GetLogRoot() string {
	return c.LogRoot
}

func GetStorageRoot() string {
	return c.StorageRoot
}

func GetDBName() string {
	return c.DBName
}

func GetDBSource() string {
	return c.DBSource
}

func GetLoggerLevel() uint8 {
	return c.LoggerLevel
}

func IsOrmLogEnabled() bool {
	return c.OrmLogEnabled
}

func GetCacheRedisAddr() string {
	return c.CacheRedisAddr
}

func GetCacheRedisPassword() string {
	return c.CacheRedisPassword
}

func GetValidateEmail() ValEmail {
	return c.ValidateEmail
}

func GetChuangLanConfig() ChuangLan {
	return c.ChuangLanConfig
}

func GetBindAddr() string {
	return c.BindAddr
}
