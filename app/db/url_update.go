package db

import (
	"time"

	"aid.im/app/util/logx"
)

func UrlUpdate(id string) {
	sql := "UPDATE `link` SET `click` =  `click` + 1 , `etime` = ? WHERE `id` = ? "
	_, err := db.Exec(sql, time.Now().Unix(), id)
	logx.Warn(err)
}
