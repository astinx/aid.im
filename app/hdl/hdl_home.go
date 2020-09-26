package hdl

import (
	"io/ioutil"
	"net/http"

	"aid.im/app/cfg"
	"aid.im/app/e"
	"aid.im/app/util/cache"
)

// home page
func Home(w http.ResponseWriter, r *http.Request) {
	key := "cache_home"
	w.WriteHeader(200)
	c := cache.Get(key)
	if c != nil {
		w.Write(c.([]byte))
		return
	}
	buf, err := ioutil.ReadFile(cfg.Opt.IndexPage)
	if err != nil {
		e.Output(w, http.StatusBadRequest, "cannot read template: "+err.Error(), nil)
		return
	}
	cache.Put(key, buf, 3600)
	w.Write(buf)
}
