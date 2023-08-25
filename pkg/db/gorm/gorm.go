package gorm

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// Gorm defines a interface for access the database.
type Gorm interface {
	DB() *gorm.DB
	SqlDB() *sql.DB
	Transaction(fc func(tx *gorm.DB) error) (err error)
	Close() error
	DropTableIfExists(value interface{}) error
}

// Config GORM Config
type Config struct {
	Debug           bool
	DBType          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	TablePrefix     string
}

// _gorm gorm struct
type _gorm struct {
	db    *gorm.DB
	sqlDB *sql.DB
}
