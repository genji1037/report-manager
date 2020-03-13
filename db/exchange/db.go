package exchange

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"report-manager/db/types"
)

var gormDb *gorm.DB

// 打开数据库
func Open(conn types.Connection) error {
	var err error
	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		conn.User, conn.Password, conn.Host, conn.Database, conn.Charset)
	gormDb, err = gorm.Open("mysql", args)
	if err != nil {
		return err
	}
	err = gormDb.DB().Ping()
	if err != nil {
		return err
	}

	gormDb.DB().SetMaxIdleConns(conn.MaxIdleConns)
	gormDb.DB().SetMaxOpenConns(conn.MaxOpenConns)

	return nil
}
