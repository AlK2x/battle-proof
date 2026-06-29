package main

import (
	"bytes"
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
	config := config.Init()
	eventBus := messaging.NewKafkaEventProducer(config)
	defer eventBus.Close()
	handler := initHttpHandler(config, eventBus)

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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	kafkaContainer, err := tckafka.Run(
		ctx,
		"confluentinc/confluent-local:7.5.0",
		tckafka.WithClusterID("test-cluster"),
		testcontainers.WithEnv(map[string]string{
			"KAFKA_NUM_PARTITIONS": "1",
		}),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = testcontainers.TerminateContainer(kafkaContainer)
	})

	brokers, err := kafkaContainer.Brokers(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, brokers)

	kafkaUrl := brokers[0]

	conn, err := kafka.Dial("tcp", kafkaUrl)
	require.NoError(t, err)
	defer conn.Close()

	err = conn.CreateTopics(
		kafka.TopicConfig{
			Topic:             messaging.JamCreatedTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		kafka.TopicConfig{
			Topic:             messaging.ParticipantAddedTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
	require.NoError(t, err)

	config := CreateTestConfig(kafkaUrl)
	eventBus := messaging.NewKafkaEventProducer(config)
	defer eventBus.Close()
	handler := initHttpHandler(config, eventBus)

	body := struct {
		Name      string    `json:"name"`
		Location  string    `json:"location"`
		Date      time.Time `json:"date"`
		CreatedBy string    `json:"created_by" binding:"uuid"`
	}{
		Name:      "John Dow",
		Location:  "KnowWhere",
		Date:      time.Date(2025, 12, 12, 0, 0, 0, 0, time.UTC),
		CreatedBy: "a181a5d7-1e56-462e-9a79-678a8e54270b",
	}
	bodyStr, err := json.Marshal(body)
	require.NoError(t, err)
	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/jams", bytes.NewBuffer(bodyStr))
	handler.ServeHTTP(r, req)

	if r.Code != http.StatusCreated {
		t.Logf("Error response body: %s", r.Body.String())
	}
	assert.Equal(t, http.StatusCreated, r.Code)

	topic := messaging.JamCreatedTopic
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
