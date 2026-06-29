//go:generate go tool oapi-codegen --config ../api/oapi-codegen.yaml ../api/openapi.yaml
package main

import (
	"context"
	"database/sql"
	"jam/config"
	"jam/internal/application"
	"jam/internal/infrastructure/messaging"
	"jam/pkg/api"
	"jam/pkg/observability"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

const AppName = "battle-proof-jam"

func openDbConnection(config config.Config) *sql.DB {
	db, err := sql.Open("mysql", config.MysqlDsn)
	if err != nil {
		log.Fatalf("error init mysql db: err %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("error mysql ping: err %v", err)
	}

	return db
}

func initHttpHandler(config config.Config, eventBus application.EventPublisher) http.Handler {
	r := gin.Default()

	service := application.NewJamService(eventBus)
	server := api.NewServer(service)
	r.Use(observability.ContextTraceMiddleware(AppName))
	r.Use(observability.LoggingMiddleware())
	api.RegisterHandlers(r, server)

	return r
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := config.Init()
	_ = openDbConnection(config)

	eventBus := messaging.NewKafkaEventProducer(config)
	defer func() {
		if err := eventBus.Close(); err != nil {
			log.Printf("error closing event bus: %v", err)
		}
	}()

	handler := initHttpHandler(config, eventBus)
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
