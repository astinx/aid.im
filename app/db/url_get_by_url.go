package db

import (
	"aid.im/app/db/model"
	"aid.im/app/util/logx"
)

func UrlGet(rawUrl string) *model.Link {
	l := &model.Link{}
	if res, err := db.Where(" url = ? ", rawUrl).Get(l); !res || err != nil {
		logx.Error(err)
		return nil
	}
	return l
}
