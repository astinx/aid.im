package db

import (
	"time"

	"aid.im/app/cfg"
	"aid.im/app/db/model"
)

// delete useless link
func UrlDel() {
	db.Where("etime < ? AND ctime < ? ", time.Now().Unix()-cfg.Opt.LinkTTl, time.Now().Unix()-cfg.Opt.LinkTTl).Delete(&model.Link{})
}
