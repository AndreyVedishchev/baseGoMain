package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"base-go/internal/utils"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// ViewCount — структура для строки в сводной таблице
type ViewCount struct {
	UUID  string
	Count int
}

// ключ, под которым в Redis хранится sorted set со счётчиками просмотров по uuid
const viewsKey = "resume:views"

type Consumer struct {
	reader *kafka.Reader
	redis  *redis.Client
	log    *zap.Logger
}

func NewConsumer(log *zap.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     utils.BuildKafkaBrokers(),
		Topic:       utils.KafkaTopic(),
		StartOffset: kafka.FirstOffset,
	})
	rdb := redis.NewClient(&redis.Options{Addr: utils.RedisAddr()})
	return &Consumer{reader: reader, redis: rdb, log: log}
}

// Run читает сообщения из topic и увеличивает счётчик по uuid в Redis
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

		if err := c.redis.ZIncrBy(context.Background(), viewsKey, 1, event.UUID).Err(); err != nil {
			c.log.Error("ошибка увеличения счётчика в Redis", zap.Error(err))
		}
	}
}

// Counts возвращает счётчики, отсортированные по убыванию количества
func (c *Consumer) Counts() []ViewCount {
	rows, err := c.redis.ZRevRangeWithScores(context.Background(), viewsKey, 0, -1).Result()
	if err != nil {
		c.log.Error("ошибка чтения счётчиков из Redis", zap.Error(err))
		return nil
	}

	result := make([]ViewCount, 0, len(rows))
	for _, z := range rows {
		result = append(result, ViewCount{UUID: z.Member.(string), Count: int(z.Score)})
	}
	return result
}

// Close закрывает консьюмер и подключение к Redis
func (c *Consumer) Close() error {
	c.redis.Close()
	return c.reader.Close()
}
