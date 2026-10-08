package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"sync"

	"base-go/internal/utils"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// ViewCount — структура для строки в сводной таблице
type ViewCount struct {
	UUID  string
	Count int
}

type Consumer struct {
	reader *kafka.Reader
	log    *zap.Logger
	mu     sync.Mutex
	counts map[string]int
}

func NewConsumer(log *zap.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     utils.BuildKafkaBrokers(),
		Topic:       utils.KafkaTopic(),
		StartOffset: kafka.FirstOffset,
	})
	return &Consumer{reader: reader, log: log, counts: make(map[string]int)}
}

// Run Читает сообщения из topic и увеличивает счётчик по uuid
func (c *Consumer) Run() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			c.log.Error("ошибка чтения сообщения из Kafka", zap.Error(err))
			continue
		}

		var event struct{ UUID string }
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("ошибка разбора сообщения из Kafka", zap.Error(err))
			continue
		}

		c.mu.Lock()
		c.counts[event.UUID]++
		c.mu.Unlock()
	}
}

// Counts возвращает текущие счётчики, отсортированные по убыванию количества
func (c *Consumer) Counts() []ViewCount {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make([]ViewCount, 0, len(c.counts))
	for uuid, count := range c.counts {
		result = append(result, ViewCount{UUID: uuid, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

// Close закрывает консьюмер
func (c *Consumer) Close() error {
	return c.reader.Close()
}
