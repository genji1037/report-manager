package types

import "report-manager/config"

// 连接信息
type Connection struct {
	Host         string
	User         string
	Password     string
	Database     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

func NewConnection(sql config.MySQL) Connection {
	return Connection{
		Host:         sql.Host,
		User:         sql.User,
		Password:     sql.Password,
		Database:     sql.Database,
		Charset:      sql.Charset,
		MaxIdleConns: sql.MaxIdleConns,
		MaxOpenConns: sql.MaxOpenConns,
	}
}
