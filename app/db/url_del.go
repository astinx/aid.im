package db

import (
	"time"

	"aid.im/app/cfg"
	"aid.im/app/db/model"
	"aid.im/app/util/logx"
)

// delete useless link
func UrlDel() {
	_, err := db.Where("etime < ? AND ctime < ? ", time.Now().Unix()-cfg.Opt.LinkTTl, time.Now().Unix()-cfg.Opt.LinkTTl).Delete(&model.Link{})
	logx.Error(err)
}
