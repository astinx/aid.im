package hdl

import (
	"net/http"

	"aid.im/app/db"
	"aid.im/app/e"
)

// del expire data
func UrlDel(w http.ResponseWriter, r *http.Request) {
	db.UrlDel()
	e.ShowNotFound(w)
}
