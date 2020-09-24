package db

import (
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var db *xorm.Engine

func InitDB(f string) {
	var err error
	db, err = xorm.NewEngine("sqlite3", f)
	if err != nil {
		panic(err)
	}
}
