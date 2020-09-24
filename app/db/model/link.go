package model

type Link struct {
	Id     string `xorm:"id" json:"id"`
	Scheme string `xorm:"scheme" json:"scheme"`
	Url    string `xorm:"url" json:"url"`
	Code   int    `xorm:"code" json:"code"`   // http status code
	Click  int64  `xorm:"click" json:"click"` // click count
	CTime  int64  `xorm:"ctime" json:"ctime"` // create time
	ETime  int64  `xorm:"etime" json:"etime"` // last click time
}
