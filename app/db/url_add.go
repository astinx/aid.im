package db

import (
	"strconv"
	"time"

	"aid.im/app/cfg"
	"aid.im/app/db/model"
	"aid.im/app/util"
	"aid.im/app/util/logx"
)

func UrlAdd(scheme, url, code, customId string, n int) *model.Link {
	if n < cfg.Opt.MinLinkLen || n > cfg.Opt.MaxLinkLen {
		n = cfg.Opt.MaxLinkLen
	}
	id := customId
	if id == "" {
		id = util.RandString(n)
	}
	if UrlGetById(id) != nil {
		return UrlAdd(scheme, url, code, "", n+1)
	}
	tempCode, _ := strconv.Atoi(code)
	l := &model.Link{
		Id:     id,
		Scheme: scheme,
		Url:    url,
		Code:   tempCode,
		Click:  0,
		CTime:  time.Now().Unix(),
		ETime:  0,
	}
	_, err := db.Insert(l)
	if err != nil {
		logx.Error(err)
		return nil
	}
	return l
}
