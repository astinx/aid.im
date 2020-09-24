package hdl

import (
	"net/http"

	"github.com/go-chi/chi"

	"aid.im/app/db"
	"aid.im/app/e"
)

func StaGet(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	if day == "" {
		day = "0"
	}
	s := db.StaGet(day)
	msg := ""
	if s == nil {
		msg = "no data"
	}
	e.Output(w, http.StatusOK, msg, s)
	return
}
