package e

import (
	"encoding/json"
	"net/http"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Output(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.WriteHeader(code)
	resp, _ := json.Marshal(Resp{code, msg, data})
	w.Write(resp)
}

func ShowNotFound(w http.ResponseWriter) {
	pageNotFound := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Not Found</title><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"><meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=0"><style>* {transition:all 0.6s;} html {height:100%;} body {font-family:Lato, sans-serif;color:#888888;margin:0;} #main {display:table;width:100%;height:100vh;text-align:center;} .fof {display:table-cell;vertical-align:middle;} .fof h1 {font-size:50px;display:inline-block;padding-right:12px;animation:type .5s alternate infinite;} @keyframes type { from {box-shadow:inset -3px 0 0 #888888;} to {box-shadow:inset -3px 0 0 transparent;} }</style></head><body><div id="main"><div class="fof"><h1>404</h1><h3>page not found</h3></div></div></body></html>`
	w.WriteHeader(404)
	w.Write([]byte(pageNotFound))
}
