package db

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var db *xorm.Engine

func InitDB(dbPath, logPath string) {
	var err error
	db, err = xorm.NewEngine("sqlite3", dbPath)
	if err != nil || db.Ping() != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(30)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	if logPath != "" {
		db.ShowSQL(true)
		db.Logger().SetLevel(0)
		var f *os.File
		if f, err = os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
			os.RemoveAll(logPath)
			if f, err = os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777); err != nil {
				panic(err)
			}
		}
		logger := log.NewSimpleLogger(f)
		logger.ShowSQL(true)
		db.SetLogger(logger)
	}
}
