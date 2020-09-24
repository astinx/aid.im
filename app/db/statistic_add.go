package db

import (
	"time"

	"aid.im/app/db/model"
	"aid.im/app/util/logx"
)

// add count for link, day = 0 is all data
func StaAdd(clickIncr, urlCountIncr int64) {
	// only accept 0 | 1
	if clickIncr < 0 || clickIncr > 1 || urlCountIncr < 0 || urlCountIncr > 1 || (clickIncr+urlCountIncr == 0) {
		return
	}
	t := time.Now()
	day := t.Year()*10000 + int(t.Month())*100 + t.Day()
	// global statistic
	exist, err := db.Where("day = 0").Get(&model.Statistics{})
	if exist == false || err != nil {
		logx.Error(err)
		db.Insert(&model.Statistics{Day: 0, TotalClick: clickIncr, TotalUrl: urlCountIncr})
	} else {
		sql := "UPDATE `statistics` SET `total_click` =  `total_click` + ? , `total_url` = `total_url`  + ?  WHERE day = 0 "
		db.Exec(sql, clickIncr, urlCountIncr)
	}
	// current day statistic
	if exist, err = db.Where("day = ? ", day).Get(&model.Statistics{}); exist == false || err != nil {
		logx.Error(err)
		db.Insert(&model.Statistics{Day: day, TotalClick: clickIncr, TotalUrl: urlCountIncr})
	} else {
		sql := "UPDATE `statistics` SET `total_click` =  `total_click` + ? , `total_url` = `total_url`  + ?  WHERE day = ? "
		db.Exec(sql, clickIncr, urlCountIncr, day)
	}

}
