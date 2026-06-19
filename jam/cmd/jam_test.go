package main

import (
	"context"
	"encoding/json"
	"jam/config"
	"jam/internal/infrastructure/messaging"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tckafka "github.com/testcontainers/testcontainers-go/modules/kafka"
)

func CreateTestConfig(kafkaUrl string) config.Config {
	return config.Config{
		MysqlDsn:  os.Getenv("MYSQL_DSN"),
		Port:      os.Getenv("APP_PORT"),
		KafkaAddr: kafkaUrl,
		Debug:     false,
	}
}

func TestSmoke(t *testing.T) {
	ctx := context.Background()
	config := config.Init()
	handler := initHttpHandler(ctx, config)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	handler.ServeHTTP(r, req)

	assert.Equal(t, http.StatusOK, r.Code)
}

type jamCreatedEvent struct {
	ID         string    `json:"event_id"`
	Type       string    `json:"event_type"`
	CreatedBy  string    `json:"created_by"`
	OccurredAt time.Time `json:"occured_at"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
}

func TestCreateJam(t *testing.T) {
	ctx := context.Background()
	// defer cancel()

	kafkaContainer, err := tckafka.Run(ctx, "confluentinc/confluent-local:7.5.0", tckafka.WithClusterID("test-cluster"))
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(kafkaContainer)
	})

	brokers, err := kafkaContainer.Brokers(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, brokers)

	kafkaUrl := brokers[0]

	topic := messaging.Topic

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     "test-group-user-created",
		Topic:       topic,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset,
		MaxWait:     500 * time.Millisecond,
	})
	defer reader.Close()

	config := CreateTestConfig(kafkaUrl)
	handler := initHttpHandler(ctx, config)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jam", nil)
	handler.ServeHTTP(r, req)

	assert.Equal(t, http.StatusCreated, r.Code)

	var got jamCreatedEvent
	require.Eventually(t, func() bool {
		readCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		msg, err := reader.ReadMessage(readCtx)
		if err != nil {
			return false
		}

		if err := json.Unmarshal(msg.Value, &got); err != nil {
			return false
		}

		return true
	}, 15*time.Second, 200*time.Millisecond)
}
