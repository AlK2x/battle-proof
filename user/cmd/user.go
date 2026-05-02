//go:generate go tool oapi-codegen --config ../api/oapi-codegen.yaml ../api/openapi.yaml
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"user/internal/app"
	"user/internal/infrastructure/mysql"
	"user/pkg/api"
	"user/pkg/api/graph"
	"user/pkg/observability"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gin-gonic/gin"
)

type GinContextKey string

const (
	AppName                  = "battle-proof-user"
	GinContext GinContextKey = "gin-context"
)

func migrateDatabase(db *sql.DB) error {
	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{
		DatabaseName:    "userdb",
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/database/mysql",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func openDbConnection(ctx context.Context, config Config) *sql.DB {
	db, err := sql.Open("mysql", config.MysqlDsn)
	if err != nil {
		log.Fatalf("error init mysql db: err %v", err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("error mysql ping: err %v", err)
	}

	return db
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContext, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinContext)
	if ginContext == nil {
		return nil, errors.New("gin.Context not found")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, errors.New("gin.Context wrong type")
	}

	return gc, nil
}

func initHttpHandler(_ context.Context, service app.UserService) http.Handler {
	r := gin.Default()
	server := api.NewServer()
	r.Use(GinContextToContextMiddleware())
	r.Use(observability.ContextTraceMiddleware(AppName))
	r.Use(observability.LoggingMiddleware())

	api.RegisterHandlers(r, server)

	r.POST("/query", grapqlHandler(service))
	r.GET("/query", grapqlHandler(service))

	return r
}

func grapqlHandler(service app.UserService) gin.HandlerFunc {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{UserService: service}}))
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	config := initConfig()
	db := openDbConnection(ctx, config)
	err := migrateDatabase(db)
	if err != nil {
		log.Fatalf("DB migration error %v", err)
	}
	userService := app.NewUserService(mysql.NewMysqlUserRepository(db))
	handler := initHttpHandler(ctx, *userService)

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
		log.Printf("shutdown server error: %v", err)
	}
}
