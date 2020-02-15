package main

import (
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
	"xorm.io/xorm"
)

var db *xorm.Engine

type Link struct {
	Id     string `xorm:"id" json:"id"`
	Scheme string `xorm:"scheme" json:"scheme"`
	Url    string `xorm:"url" json:"url"`
	Code   int    `xorm:"code" json:"code"`   // http status code
	Click  int64  `xorm:"click" json:"click"` // click count
	CTime  int64  `xorm:"ctime" json:"ctime"` // create time
	ETime  int64  `xorm:"etime" json:"etime"` // last click time
}

type Statistics struct {
	Day        int   `xorm:"day" json:"day"`
	TotalClick int64 `xorm:"total_click" json:"total_click"`
	TotalUrl   int64 `xorm:"total_url" json:"total_url"`
}

func NewDB(f string) {
	var err error
	db, err = xorm.NewEngine("sqlite3", f)
	if err != nil {
		panic(err)
	}
}

// url
func dbUrlGetById(id string) *Link {
	var l Link
	if res, _ := db.Where("id = ?", id).Get(&l); res == false {
		return nil
	}
	return &l
}

func dbUrlGet(scheme, rawUrl string) *Link {
	l := &Link{}
	if res, err := db.Where(" scheme = ? AND url = ? ", scheme, rawUrl).Get(l); !res || err != nil {
		return nil
	}
	return l
}

// delete useless link
func dbUrlDel() {
	db.Where("etime < ? AND ctime < ? ", time.Now().Unix()-opts.LinkTimeout, time.Now().Unix()-opts.LinkTimeout).Delete(&Link{})
}

func dbUrlAdd(scheme, url, code string, n int) *Link {
	if n < opts.MinLinkLen || n > opts.MaxLinkLen {
		n = opts.MaxLinkLen
	}
	id := RandString(n)
	if dbUrlGetById(id) != nil {
		return dbUrlAdd(scheme, url, code, n+1)
	}
	tempCode, _ := strconv.Atoi(code)
	l := &Link{
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
		return nil
	}
	return l
}

func dbUrlUpdate(id string) {
	sql := "UPDATE `link` SET `click` =  `click` + 1 , `etime` = ? WHERE `id` = ? "
	db.Exec(sql, time.Now().Unix(), id)
}

// add count for link, day = 0 is all data
func dbStaAdd(clickIncr, urlCountIncr int64) {
	// only accept 0 | 1
	if clickIncr < 0 || clickIncr > 1 || urlCountIncr < 0 || urlCountIncr > 1 || (clickIncr+urlCountIncr == 0) {
		return
	}
	t := time.Now()
	day := t.Year()*10000 + int(t.Month())*100 + t.Day()
	// global statistic
	exist, err := db.Where("day = 0").Get(&Statistics{})
	if exist == false || err != nil {
		db.Insert(&Statistics{Day: 0, TotalClick: clickIncr, TotalUrl: urlCountIncr})
	} else {
		sql := "UPDATE `statistics` SET `total_click` =  `total_click` + ? , `total_url` = `total_url`  + ?  WHERE day = 0 "
		db.Exec(sql, clickIncr, urlCountIncr)
	}
	// current day statistic
	if exist, err = db.Where("day = ? ", day).Get(&Statistics{}); exist == false || err != nil {
		db.Insert(&Statistics{Day: day, TotalClick: clickIncr, TotalUrl: urlCountIncr})
	} else {
		sql := "UPDATE `statistics` SET `total_click` =  `total_click` + ? , `total_url` = `total_url`  + ?  WHERE day = ? "
		db.Exec(sql, clickIncr, urlCountIncr, day)
	}

}

// get statistics
func dbStaGet(day string) *Statistics {
	var s Statistics
	if res, _ := db.Where("day = ? ", day).Get(&s); res == false {
		return nil
	}
	return &s
}
