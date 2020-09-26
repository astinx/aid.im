package db

import (
	"aid.im/app/db/model"
	"aid.im/app/util/logx"
)

// get statistics
func StaGet(day string) *model.Statistics {
	var s model.Statistics
	if res, err := db.Where("day = ? ", day).Get(&s); err != nil || res == false {
		logx.Error(err)
		return nil
	}
	return &s
}
