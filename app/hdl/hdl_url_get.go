package hdl

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi"

	"aid.im/app/db"
	"aid.im/app/db/model"
	"aid.im/app/e"
	"aid.im/app/util"
	"aid.im/app/util/cache"
)

func UrlGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		e.ShowNotFound(w)
		return
	}
	key := "cache_url_" + id
	var l *model.Link
	if tmp := cache.Get(key); tmp == nil {
		l = db.UrlGetById(id)
		if l == nil {
			e.ShowNotFound(w)
			return
		}
		l.Url, _ = url.QueryUnescape(l.Url)
		cache.Put(key, l, 600)
	} else {
		l = tmp.(*model.Link)
	}
	// update statistics
	go db.UrlUpdate(id)
	go db.StaAdd(1, 0)
	if l.Code == 200 {
		res := "<!DOCTYPE html><html><head></head><body><script>window.location.href= '" + l.Scheme + "://" + l.Url + "';</script></body></html>"
		if !util.IsURL(l.Scheme + "://" + l.Url) {
			res = "<!DOCTYPE html><html><head><meta name='viewport'content='maximum-scale=1.0,minimum-scale=1.0,user-scalable=0,width=device-width,initial-scale=1.0'><style>body{display:-ms-flexbox;display:-webkit-box;display:flex;-ms-flex-align:center;-ms-flex-pack:center;-webkit-box-align:center;align-items:center;-webkit-box-pack:center;justify-content:center;font-size:2rem;padding-top:30px}</style></head><body><div id='box'></div><script>if(/Android|webOS|iPhone|iPod|BlackBerry/i.test(navigator.userAgent)){window.location.href='" + l.Scheme + "://" + l.Url + "'}else{document.getElementById('box').innerHTML='Please open in mobile device'}</script></body></html>"
		}
		w.Write([]byte(res))
		return
	}
	http.Redirect(w, r, l.Scheme+"://"+l.Url, l.Code)
}
