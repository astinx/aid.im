package hdl

import (
	"io/ioutil"
	"net/http"

	"aid.im/app/cfg"
	"aid.im/app/util/cache"

	"aid.im/app/e"
)

// home page
func Home(w http.ResponseWriter, r *http.Request) {
	key := "tinyurl_home"
	w.WriteHeader(200)
	c := cache.Get(key)
	if c != nil {
		w.Write(c.([]byte))
		return
	}
	f, err := ioutil.ReadFile(cfg.Opt.IndexPage)
	if err != nil {
		e.Output(w, http.StatusBadRequest, "cannot read template: "+err.Error(), nil)
		return
	}
	cache.Put(key, f, 3600)
	w.Write(f)
}
