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
		ShowNotFound(w)
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
	resp(w, http.StatusOK, msg, s)
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
	addCors(w)
	s := SplitUrl(rawUrl)
	if s == nil {
		resp(w, http.StatusBadRequest, "invalid url", nil)
		return
	}
	if !IsSupportType(rawCode) {
		resp(w, http.StatusBadRequest, "nonsupport type code", nil)
		return
	}
	if IsURL(rawUrl) {
		//is same with current host
		l, err := url.Parse(rawUrl)
		if err != nil || l.Host == opts.Host {
			resp(w, http.StatusBadRequest, "invalid url", nil)
			return
		}
	} else if IsApp(rawUrl) { // is app link
		rawCode = "200"
	} else {
		resp(w, http.StatusBadRequest, "invalid url", nil)
		return
	}
	customId := r.Form.Get("id")
	if customId != "" {
		if len(customId) < 4 || len(customId) > 12 {
			resp(w, http.StatusBadRequest, "custom id length must between 4 - 12 ", nil)
			return
		}
		// not use regex cause loop faster than regex
		for _, v := range customId {
			if v < 'a' && v > 'z' && v < 'A' && v > 'Z' && v != '_' {
				resp(w, http.StatusBadRequest, "custom id characters only accept 0-9,A-Z,a-z AND _", nil)
				return
			}
		}
		if dbUrlGetById(customId) != nil {
			resp(w, http.StatusBadRequest, "custom id exist", nil)
			return
		}
	}
	l := dbUrlGet(s[1])
	if l != nil && (s[0] == "http" || s[0] == "https") && (l.Scheme == "http" || l.Scheme == "https") {
		resp(w, http.StatusOK, "", l)
		return
	}
	res := dbUrlAdd(s[0], s[1], rawCode, customId, opts.MinLinkLen)
	if res == nil {
		resp(w, http.StatusBadRequest, "unknown error", nil)
		return
	}
	go dbStaAdd(0, 1)
	resp(w, http.StatusOK, "", res)
	return
}

func addCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                                                                                                                                  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,Content-Length,X-CSRF-Token,Accept,X-Requested-With, Authorization, Token,Origin, No-Cache, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, X-E4M-With") //header的类型
	w.Header().Add("Access-Control-Allow-Credentials", "true")                                                                                                                                                                                          //设置为true，允许ajax异步请求带cookie信息
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                                                                                                                                                                   //允许请求方法
	w.Header().Set("content-type", "application/json;charset=UTF-8")                                                                                                                                                                                    //返回数据格式是json
}
