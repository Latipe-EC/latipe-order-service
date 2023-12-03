package gorm

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

// SqlDB returns `*sql.DB`
func (g *_gorm) SqlDB() *sql.DB {
	return g.sqlDB
}

func (g *_gorm) Exec(fc func(tx *gorm.DB) error, ctx context.Context) (err error) {
	panicked := true
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := g.db.WithContext(ctx).Begin()
	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	err = fc(tx)

	if err == nil {
		err = tx.Commit().Error
	}

	panicked = false
	return
}

// DB returns `*gorm.DB`
func (g *_gorm) DB() *gorm.DB {
	return g.db
}

// Close close current db connection. If database connection is not an io.Closer, returns an error.
func (g *_gorm) Close() error {
	return g.sqlDB.Close()
}

// DropTableIfExists drop table if it is exist
func (g *_gorm) DropTableIfExists(value interface{}) error {
	return g.db.Migrator().DropTable(value)
}

// Transaction start a transaction as a block.
// If it is failed, will rollback and return error.
// If it is successful, will commit.
// ref: https://github.com/jinzhu/gorm/blob/master/main.go#L533
func (g *_gorm) Transaction(fc func(tx *gorm.DB) error) (err error) {
	panicked := true
	tx := g.db.Begin()
	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	err = fc(tx)

	if err == nil {
		err = tx.Commit().Error
	}

	panicked = false
	return
}
