package db

import (
	"time"
)

func UrlUpdate(id string) {
	sql := "UPDATE `link` SET `click` =  `click` + 1 , `etime` = ? WHERE `id` = ? "
	db.Exec(sql, time.Now().Unix(), id)
}
