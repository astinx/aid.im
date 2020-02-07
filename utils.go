package main

import (
	crand "crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	maxURLRuneCount = 2083
	minURLRuneCount = 3

	IP           string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLSchema    string = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername  string = `(\S+(:\S*)?@)`
	URLPath      string = `((\/|\?|#)[^\s]*)`
	URLPort      string = `(:(\d{1,5}))`
	URLIP        string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	URLSubdomain string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URL          string = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	URL_APP      string = `^[a-z]{2,10}:\/\/\w`
)

// IsURL check if the string is an URL.
func IsURL(str string) bool {
	if str == "" || utf8.RuneCountInString(str) >= maxURLRuneCount || len(str) <= minURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	strTemp := str
	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
		// support no indicated urlscheme but with colon for port number
		// http:// is appended so url.Parse will succeed, strTemp used so it does not impact rxURL.MatchString
		strTemp = "http://" + str
	}
	u, err := url.Parse(strTemp)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return regexp.MustCompile(URL).MatchString(str)
}

func IsApp(str string) bool {
	return regexp.MustCompile(URL_APP).MatchString(str)
}

func IsSupportType(s string) bool {
	if len(s) != 3 {
		return false
	}
	supportType := []string{"200", "301", "302", "303", "307"}
	for _, v := range supportType {
		if s == v {
			return true
		}
	}
	return false
}

func SplitUrl(str string) []string {
	s := strings.SplitN(str, "://", 2)
	if len(s) < 2 {
		return nil
	}
	s[0] = strings.ToLower(s[0])
	return s
}

func ShowNotFound(w http.ResponseWriter) {
	pageNotFound := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Not Found</title><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=0"><style>* {transition:all 0.6s;} html {height:100%;} body {font-family:Lato, sans-serif;color:#888888;margin:0;} #main {display:table;width:100%;height:100vh;text-align:center;} .fof {display:table-cell;vertical-align:middle;} .fof h1 {font-size:50px;display:inline-block;padding-right:12px;animation:type .5s alternate infinite;} @keyframes type { from {box-shadow:inset -3px 0 0 #888888;} to {box-shadow:inset -3px 0 0 transparent;} }</style></head><body><div id="main"><div class="fof"><h1>404</h1><h3>page not found</h3></div></div></body></html>`
	w.WriteHeader(404)
	w.Write([]byte(pageNotFound))
}

func RandString(l int) string {
	var (
		str string
		num *big.Int
		chr int64
	)
	for len(str) < l {
		num, _ = crand.Int(crand.Reader, big.NewInt(123))
		chr = num.Int64()
		if (chr >= 48 && chr <= 57) || (chr >= 65 && chr <= 90) || (chr >= 97 && chr <= 122) {
			str += string(chr)
		}
	}
	return str
}

func resp(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.WriteHeader(code)
	resp, _ := json.Marshal(Resp{code, msg, data})
	w.Write(resp)
}
