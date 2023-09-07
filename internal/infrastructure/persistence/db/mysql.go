package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	gormDB "gorm.io/gorm"
	"order-service-rest-api/config"
	"order-service-rest-api/pkg/db/gorm"
)

func NewMySQLConnection(configuration *config.Config) *gormDB.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		configuration.DB.Mysql.UserName,
		configuration.DB.Mysql.Password,
		configuration.DB.Mysql.Host,
		configuration.DB.Mysql.Port,
		configuration.DB.Mysql.Database,
	)
	cfg := gorm.Config{
		DSN:             dataSourceName,
		MaxOpenConns:    configuration.DB.Mysql.MaxOpenConns,
		MaxIdleConns:    configuration.DB.Mysql.MaxIdleConns,
		ConnMaxLifetime: configuration.DB.Mysql.ConnMaxLifetime,
		ConnMaxIdleTime: configuration.DB.Mysql.ConnMaxIdleTime,
		DBType:          "mysql",
	}
	dial := mysql.Open(cfg.DSN)
	conn, err := gormDB.Open(dial)
	if err != nil {
		fmt.Printf("ERORR:%v", err.Error())
		return nil
	}

	return conn
}
