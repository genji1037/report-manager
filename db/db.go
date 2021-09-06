package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
	"report-manager/db/types"
	"report-manager/logger"
	"time"
)

var gormDb *gorm.DB

type Table interface {
	Create() error
}

var tables = []Table{
	new(SieCount),
	new(ExchangeSpecialUser),
	new(ExchangeSpecialUserReport),
}

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
	gormDb.DB().SetConnMaxLifetime(7 * time.Hour)
	for i := 0; i < len(tables); i++ {
		ensureTable(tables[i])
	}
	return nil
}

// 确保表存在
func ensureTable(table interface{}) {
	typ := reflect.TypeOf(table)
	tablename := typ.String()[1:]
	if !gormDb.HasTable(table) {
		logger.Infof("Creating table: %s", tablename)
		if err := gormDb.CreateTable(table).Error; err != nil {
			logger.Warnf("Failed to create table %s, %v", tablename, err)
		}
	} else { // always migrate
		logger.Infof("Auto migrate table: %s", tablename)
		if err := gormDb.AutoMigrate(table).Error; err != nil {
			logger.Warnf("Failed to auto migrate table %s, %v", tablename, err)
		}
	}
}
