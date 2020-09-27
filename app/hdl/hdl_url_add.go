package hdl

import (
	"net/http"
	"net/url"
	"strings"

	"aid.im/app/cfg"
	"aid.im/app/db"
	"aid.im/app/e"
	"aid.im/app/util"
)

// add a new one link
func UrlAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rawUrl := strings.Trim(r.Form.Get("url"), "/")
	rawCode := r.Form.Get("code")
	if rawCode == "" {
		rawCode = "200"
	}
	s := util.SplitUrl(rawUrl)
	if s == nil {
		e.Output(w, http.StatusBadRequest, "invalid url, split err", nil)
		return
	}
	if !util.IsSupportType(rawCode) {
		e.Output(w, http.StatusBadRequest, "nonsupport type code", nil)
		return
	}
	if util.IsURL(rawUrl) {
		//is same with current host
		l, err := url.Parse(rawUrl)
		if err != nil || l.Host == r.Host {
			e.Output(w, http.StatusBadRequest, "invalid url: can't use "+r.Host, nil)
			return
		}
	} else if util.IsApp(rawUrl) { // is app link
		rawCode = "200"
	} else {
		e.Output(w, http.StatusBadRequest, "invalid url", nil)
		return
	}
	customId := r.Form.Get("id")
	if customId != "" {
		if len(customId) < 4 || len(customId) > 12 {
			e.Output(w, http.StatusBadRequest, "custom id length must between 4 - 12 ", nil)
			return
		}
		// not use regex cause loop faster than regex
		for _, v := range customId {
			if v < 'a' && v > 'z' && v < 'A' && v > 'Z' && v != '_' {
				e.Output(w, http.StatusBadRequest, "custom id characters only accept 0-9,A-Z,a-z AND _", nil)
				return
			}
		}
		if db.UrlGetById(customId) != nil {
			e.Output(w, http.StatusBadRequest, "custom id exist", nil)
			return
		}
	}
	l := db.UrlGet(s[1])
	if l != nil && l.Url == url.QueryEscape(s[1]) {
		l.Url, _ = url.QueryUnescape(l.Url)
		e.Output(w, http.StatusOK, "", l)
		return
	}
	res, err := db.UrlAdd(s[0], url.QueryEscape(s[1]), rawCode, customId, cfg.Opt.MinLinkLen)
	if err != nil {
		e.Output(w, http.StatusBadRequest, "err : "+err.Error(), nil)
		return
	}
	go db.StaAdd(0, 1)
	e.Output(w, http.StatusOK, "", res)
	return
}
