package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSmoke(t *testing.T) {
	ctx := context.Background()
	handler := initHttpHandler(ctx, nil)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	handler.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
}

func TestBattleService(t *testing.T) {
	dsn, err := runTestDbContainers()
	if err != nil {
		t.Fatalf("runTestContainers error %v", err)
	}

	config := Config{
		MysqlDsn:      dsn,
		Port:          "8080",
		MigrationPath: "file://../database/mysql",
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
	}

	handler, err := createHandler(context.Background(), config)
	assert.NoError(t, err)
	srv := httptest.NewServer(handler)

	req, err := http.NewRequest("GET", srv.URL+"/health", nil)
	assert.NoError(t, err)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func runTestDbContainers() (string, error) {
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
		return "", err
	}
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return "", err
	}

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		"testuser", "testpass", host, port.Port(), "testdb")
	return dns, nil
}
