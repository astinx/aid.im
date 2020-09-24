package db

import (
	"aid.im/app/db/model"
)

// url
func UrlGetById(id string) *model.Link {
	var l model.Link
	if res, _ := db.Where("id = ?", id).Get(&l); res == false {
		return nil
	}
	return &l
}
