package db

import (
	"aid.im/app/db/model"
)

// get statistics
func StaGet(day string) *model.Statistics {
	var s model.Statistics
	if res, _ := db.Where("day = ? ", day).Get(&s); res == false {
		return nil
	}
	return &s
}
