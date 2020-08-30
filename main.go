package main

import (
	"flag"
	"github.com/go-chi/chi"
	"net/http"
	"os"
)

func main() {
	cfgPath := flag.String("opt", "app.yml", "config file path")
	flag.Parse()
	var err error
	if _, err = os.Stat(*cfgPath); err != nil {
		panic(err.Error())
	}
	NewCfg(*cfgPath)
	NewDB(opts.DB)
	NewCache()

	r := chi.NewRouter()
	r.Get("/", hdlHome)
	r.Get("/api", hdlAdd)
	r.Post("/api", hdlAdd)
	r.Get("/del", hdlDel)
	r.Get("/s/{day:[0-9]+}", hdlStatistic)
	r.Get("/{id:[a-zA-Z0-9]+}", hdlGet)
	http.ListenAndServe(":8080", r)
}
