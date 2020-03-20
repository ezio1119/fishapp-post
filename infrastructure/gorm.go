package infrastructure

import (
	"log"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func NewGormDB() *gorm.DB {
	mysqlConf := &mysql.Config{
		User:                 conf.C.Db.User,
		Passwd:               conf.C.Db.Pass,
		Net:                  conf.C.Db.Net,
		Addr:                 conf.C.Db.Host + ":" + conf.C.Db.Port,
		DBName:               conf.C.Db.Name,
		ParseTime:            conf.C.Db.Parsetime,
		Loc:                  time.Local,
		AllowNativePasswords: conf.C.Db.AllowNativePasswords,
	}

	db, err := gorm.Open(conf.C.Db.Dbms, mysqlConf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	if conf.C.Sv.Debug {
		db.LogMode(true)
	}
	return db
}
