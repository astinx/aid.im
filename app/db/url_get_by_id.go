package db

import (
	"aid.im/app/db/model"
	"aid.im/app/util/logx"
)

// url
func UrlGetById(id string) *model.Link {
	var l model.Link
	if res, err := db.Where("id = ?", id).Get(&l); res == false || err != nil {
		logx.Warn(err)
		return nil
	}
	return &l
}
