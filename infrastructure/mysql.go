package infrastructure

import (
	"database/sql"
	"log"
	"time"

	"github.com/ezio1119/fishapp-post/conf"
	"github.com/go-sql-driver/mysql"
)

func NewMysqlConn() *sql.DB {
	mysqlConf := &mysql.Config{
		User:                 conf.C.Db.User,
		Passwd:               conf.C.Db.Pass,
		Net:                  conf.C.Db.Net,
		Addr:                 conf.C.Db.Host + ":" + conf.C.Db.Port,
		DBName:               conf.C.Db.Name,
		ParseTime:            conf.C.Db.Parsetime,
		Loc:                  time.Now().Location(),
		AllowNativePasswords: conf.C.Db.AllowNativePasswords,
	}

	conn, err := sql.Open(conf.C.Db.Dbms, mysqlConf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err = conn.Ping(); err != nil {
		log.Fatal(err)
	}
	return conn
}
