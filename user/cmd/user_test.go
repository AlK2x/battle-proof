package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user/internal/app"
	"user/internal/infrastructure/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSmoke(t *testing.T) {
	ctx := context.Background()
	userService := app.NewUserService(mysql.NewMysqlUserRepository(nil))
	handler := initHttpHandler(ctx, *userService)

	r := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(ctx, "GET", "/health", nil)
	handler.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
}

type GrapQLRequest struct {
	Query     string
	Variables map[string]any
}

func TestUserService(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "test",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("ready for connection"),
			wait.ForListeningPort("3306/tcp"),
		).WithDeadline(1 * time.Minute),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start testcontainer error %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("get testcontainer host error %v", err)
	}
	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("get testcontainer port error %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		"testuser", "testpass", host, port.Port(), "testdb")

	config := Config{
		MysqlDsn:      dsn,
		Port:          "8080",
		MigrationPath: "file://../database/mysql",
	}

	db, err := openDbConnection(ctx, config)
	if err != nil {
		t.Fatalf("open DB connection error %v", err)
	}
	err = migrateDatabase(db, config)
	if err != nil {
		t.Fatalf("DB migration error %v", err)
	}
	userService := app.NewUserService(mysql.NewMysqlUserRepository(db))

	handler := initHttpHandler(ctx, *userService)

	srv := httptest.NewServer(handler)

	query := "mutation($input: CreateUserInput!) { createUser(input: $input) { id name email level createdAt } }"
	variables := map[string]any{
		"input": map[string]any{
			"name":  "Jane Dow",
			"email": "jane.dow@local.test",
			"level": "pro",
		},
	}
	executeQuery(t, srv.URL+"/query", query, variables)
}

func executeQuery(t *testing.T, url string, query string, variables map[string]any) {
	gqlRequest := &GrapQLRequest{
		Query:     query,
		Variables: variables,
	}
	body, err := json.Marshal(gqlRequest)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
