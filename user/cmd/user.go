package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"user/pkg/observability"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func openDbConnection(config Config) *sql.DB {
	db, err := sql.Open("mysql", config.MysqlDsn)
	if err != nil {
		log.Fatalf("error init mysql db: err %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error mysql ping: err %v", err)
	}

	return db
}

func initHttpHandler(_ context.Context, _ *sql.DB) http.Handler {
	r := gin.Default()
	r.Use(observability.ContextTraceMiddleware())
	r.GET("/health", func(ctx *gin.Context) {})

	return r
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	config := initConfig()
	db := openDbConnection(config)
	handler := initHttpHandler(ctx, db)

	server := &http.Server{
		Addr:         config.Port,
		Handler:      handler,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	stop()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
