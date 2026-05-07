//go:generate go tool oapi-codegen --config ../api/oapi-codegen.yaml ../api/openapi.yaml
package main

import (
	"battle/internal/infrastructure/mysql"
	"battle/pkg/api"
	"battle/pkg/observability"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gin-gonic/gin"
)

const AppName = "battle-proof-jam"

func openDbConnection(config Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.MysqlDsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql db error %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql db error %v", err)
	}

	return db, nil
}

func migrateDatabase(db *sql.DB, config Config) error {
	driver, err := mysqlmigrate.WithInstance(db, &mysqlmigrate.Config{
		DatabaseName:    "userdb",
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationPath,
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

func initHttpHandler(_ context.Context, db *sql.DB) http.Handler {
	r := gin.Default()

	battleRepo := mysql.NewMysqlBattleRepository(db)
	scoreRepo := mysql.NewMysqlScoreRepository(db)
	battleService := CreateBattleService(battleRepo, scoreRepo)
	server := api.NewServer(battleService)
	r.Use(observability.ContextTraceMiddleware(AppName))
	r.Use(observability.LoggingMiddleware())
	api.RegisterHandlers(r, server)

	return r
}

func createHandler(ctx context.Context, config Config) (http.Handler, error) {
	db, err := openDbConnection(config)
	if err != nil {
		return nil, err
	}
	err = migrateDatabase(db, config)
	if err != nil {
		return nil, err
	}
	handler := initHttpHandler(ctx, db)
	return handler, nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	config := initConfig()
	handler, err := createHandler(ctx, config)
	if err != nil {
		log.Fatalf("createHandlerError %v", err)
	}

	server := &http.Server{
		Addr:         config.Port,
		Handler:      handler,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Printf("running service on %v", config.Port)

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
