package model

type Statistics struct {
	Day        int   `xorm:"day" json:"day"`
	TotalClick int64 `xorm:"total_click" json:"total_click"`
	TotalUrl   int64 `xorm:"total_url" json:"total_url"`
}
