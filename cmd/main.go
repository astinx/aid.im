package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"

	"aid.im/app/cfg"
	"aid.im/app/db"
	"aid.im/app/hdl"
	"aid.im/app/mw"
)

func init() {
	cfg.InitCfg()
	cfg.InitCache()
	cfg.InitLog()
	db.InitDB(cfg.DBPath)
}
func initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(mw.Cors)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		return
	})
	return r
}

func main() {
	r := initRouter()
	route(r)
	srv := &http.Server{
		Addr:    cfg.SrvAddr,
		Handler: r,
	}
	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		fmt.Println("\nListening and serving HTTP on " + cfg.SrvAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}

func route(r *chi.Mux) {
	r.Get("/", hdl.Home)
	r.Get("/api", hdl.UrlAdd)
	r.Post("/api", hdl.UrlAdd)
	r.Get("/del", hdl.UrlDel)
	r.Get("/s/{day:[0-9]+}", hdl.StaGet)
	r.Get("/{id:[a-zA-Z0-9]+}", hdl.UrlGet)
}
