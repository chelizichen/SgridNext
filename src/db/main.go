package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
)

// DB 全局数据库实例
var DB *gorm.DB

func initDB(dsn string, dbType string) (*gorm.DB, error) {
	fmt.Printf("init db %s | %s", dsn, dbType)
	var err error
	if dbType == constant.DB_TYPE_MYSQL {
		fmt.Printf("INIT DB WITH MYSQL")
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect database: %v", err)
		}
	}
	if dbType == constant.DB_TYPE_POSTGRES {
		fmt.Printf("INIT DB WITH POSTGRES")
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect database: %v", err)
		}
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %v", err)
	}
	// 连接池配置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return DB, nil

}

// defaultInitDB 默认(mysql)初始化数据库连接
func defaultInitDB() (*gorm.DB, error) {
	fmt.Println("get conf ", config.Conf)
	dsn := config.Conf.Get("db")
	dbType := config.Conf.Get("dbtype")
	return initDB(dsn, dbType)
}

// InitDB 初始化数据库连接
func InitDB(args ...string) (*gorm.DB, error) {
	if len(args) == 0 {
		return defaultInitDB()
	}
	// 根据传入的参数选择数据库类型
	var dsn string = args[0]
	var dbType string = args[1]
	if dsn == "" || dbType == "" {
		return nil, fmt.Errorf("dsn and dbType must be provided")
	}
	return initDB(dsn, dbType)
}
