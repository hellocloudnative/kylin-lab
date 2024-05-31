package database

import "kylin-lab/tools/config"

var (
	DbType   string
	Host     string
	Port     int
	Name     string
	Username string
	Password string
)

func Setup() {
	dbType := config.DatabaseConfig.Dbtype
	if dbType == "mysql" {
		var db = new(Mysql)
		db.Setup()
	}
}
