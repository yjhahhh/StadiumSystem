package connection

import (
	"fmt"
	"sync"
	"time"

	"common/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var one sync.Once

func GetDB() *gorm.DB {
	return db
}

// const dsn = "root:system@tcp(192.168.201.1:3306)/stadium?charset=utf8&parseTime=True&loc=Local"

func InitMySQL() {
	one.Do(func() {
		err := newConnection()
		if err != nil {
			panic(err)
		}
	})
}

func newConnection() error {
	conf := config.GetMySQLConfig()
	if conf == nil {
		return fmt.Errorf("MySQLConfig is nil")
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=%s&parseTime=True&loc=Local", conf.UserName, conf.Password, conf.Protocol, conf.Address, conf.DBName, conf.Charset)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	db = conn
	
	return nil
}