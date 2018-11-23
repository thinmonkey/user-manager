package db

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/thinmonkey/user-manager/utils/log"
)

var dbEngine *gorm.DB

//配置数据库引擎
func ConfigDbEngine(host string, port string, user string, password string, dbName string) {
	dbArgs := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true`,
		user, password,
		host, port, dbName)
	log.Info("mysql_dbconnectconfig:" + dbArgs)
	var err error
	dbEngine, err = gorm.Open("mysql", dbArgs)
	if err != nil {
		panic(err)
	}
}

func Db() *gorm.DB {
	return dbEngine
}
