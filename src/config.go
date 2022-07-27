package main

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

func GetMysqlConfig() *mysql.Config {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = "127.0.0.1:3306"
	mysqlConfig.DBName = "test"
	mysqlConfig.User = "root"
	mysqlConfig.Passwd = "root"
	mysqlConfig.Timeout = 30 * time.Second
	mysqlConfig.ReadTimeout = 1 * time.Second
	mysqlConfig.WriteTimeout = 1 * time.Second
	//如果Table有用到datetime，要加上parseTime=True，不然解析不了這個Type。
	mysqlConfig.ParseTime = true

	return mysqlConfig
}
