package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
	"strings"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// home page
func hdlHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
}

// get link
func hdlGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {

	}
	l := dbUrlGetById(id)
	if l == nil {
		ShowNotFound(w)
		return
	}
	// update statistics
	go dbUrlUpdate(id)
	go dbStaAdd(1, 0)
	if l.Code == 200 {
		if !IsURL(l.Scheme + "://" + l.Url) {
			w.Write([]byte("<!DOCTYPE html><html><head><meta name='viewport'content='maximum-scale=1.0,minimum-scale=1.0,user-scalable=0,width=device-width,initial-scale=1.0'><style>body{display:-ms-flexbox;display:-webkit-box;display:flex;-ms-flex-align:center;-ms-flex-pack:center;-webkit-box-align:center;align-items:center;-webkit-box-pack:center;justify-content:center;font-size:2rem;padding-top:30px}</style></head><body><div id='box'></div><script>if(/Android|webOS|iPhone|iPod|BlackBerry/i.test(navigator.userAgent)){window.location.href='" + l.Scheme + "://" + l.Url + "'}else{document.getElementById('box').innerHTML='Please open in mobile device'}</script></body></html>"))
			return
		}
		w.Write([]byte("<!DOCTYPE html><html><head></head><body><script>window.location.href= '" + l.Scheme + "://" + l.Url + "';</script></body></html>"))
		return
	}
	http.Redirect(w, r, l.Scheme+"://"+l.Url, l.Code)
}

// del expire data
func hdlDel(w http.ResponseWriter, r *http.Request) {
	dbUrlDel()
	ShowNotFound(w)
}

func hdlStatistic(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	if day == "" {
		day = "0"
	}
	s := dbStaGet(day)
	msg := ""
	if s == nil {
		msg = "no data"
	}
	resp(w, 200, msg, s)
	return
}

// add a new one link
func hdlAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rawUrl := strings.Trim(r.Form.Get("url"), "/")
	rawCode := r.Form.Get("code")
	if rawCode == "" {
		rawCode = "200"
	}
	s := SplitUrl(rawUrl)
	if s == nil {
		resp(w, 400, "invalid url", nil)
		return
	}
	if !IsSupportType(rawCode) {
		resp(w, 400, "nonsupport type code", nil)
		return
	}
	if IsURL(rawUrl) {
		//is same with current host
		l, err := url.Parse(rawUrl)
		if err != nil || l.Host == opts.Host {
			resp(w, 400, "invalid url", nil)
			return
		}
	} else if IsApp(rawUrl) { // is app link
		rawCode = "200"
	} else {
		resp(w, 400, "invalid url", nil)
		return
	}
	l := dbUrlGet(s[0], s[1])
	if l != nil {
		resp(w, 200, "", l)
		return
	}
	res := dbUrlAdd(s[0], s[1], rawCode, opts.MinLinkLen)
	if res == nil {
		resp(w, 400, "unknown error", nil)
		return
	}
	go dbStaAdd(0, 1)
	resp(w, 200, "", res)
	return
}
