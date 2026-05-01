package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"user/internal/app"
	"user/internal/infrastructure/mysql"

	"github.com/stretchr/testify/assert"
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
