//go:generate go tool oapi-codegen --config ../api/oapi-codegen.yaml ../api/openapi.yaml
package main

import (
	"context"
	"gradebook/pkg/api"
	"gradebook/pkg/observability"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/gin-gonic/gin"
)

const AppName = "battle-proof-user"

func openDbConnection(ctx context.Context, config Config) *mongo.Client {
	mongo, err := mongo.Connect(options.Client().ApplyURI(config.MongoDSN))
	if err != nil {
		log.Fatalf("error mongo connect: err %v dsn: %s", err, config.MongoDSN)
	}
	err = mongo.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error mongo ping: err %v", err)
	}
	return mongo
}

func initHttpHandler(_ context.Context) http.Handler {
	r := gin.Default()
	server := api.NewServer()
	r.Use(observability.ContextTraceMiddleware(AppName))
	r.Use(observability.LoggingMiddleware())
	api.RegisterHandlers(r, server)

	return r
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	config := initConfig()
	_ = openDbConnection(ctx, config)
	handler := initHttpHandler(ctx)

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
