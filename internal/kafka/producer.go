package kafka

import (
	"base-go/internal/utils"
	"context"
	"encoding/json"
	"time"

	"base-go/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Producer публикует события о резюме в Kafka
type Producer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

// NewProducer создаёт продюсер, настроенный через переменные окружения KAFKA_BROKERS / KAFKA_TOPIC
func NewProducer(log *zap.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(utils.BuildKafkaBrokers()...),
		Topic: utils.KafkaTopic(),
		// без явного Balancer партиция выбиралась бы round-robin'ом и ключ (uuid)
		// на неё бы не влиял — тогда порядок сообщений по одному резюме не гарантирован
		Balancer: &kafka.Hash{},
		// в dev-окружении topic может быть ещё не создан — позволяем брокеру создать его на лету
		AllowAutoTopicCreation: true,
		// защита на уровне Writer, если брокер не отвечает; независима от ctx, переданного в Publish
		WriteTimeout: 5 * time.Second,
	}
	return &Producer{writer: writer, log: log}
}

// Publish сериализует резюме в JSON и публикует в Kafka с ключом uuid
func (p *Producer) Publish(resume *models.Resume) error {
	value, err := json.Marshal(resume)
	if err != nil {
		p.log.Error("ошибка сериализации резюме для Kafka", zap.String("uuid", resume.UUID), zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Key:   []byte(resume.UUID),
		Value: value,
	}

	// WriteMessages блокируется до подтверждения от брокера; context.Background() — без
	// отмены и дедлайна извне, таймаут ожидания задаёт только WriteTimeout у Writer
	if err := p.writer.WriteMessages(context.Background(), msg); err != nil {
		p.log.Error("ошибка публикации сообщения в Kafka", zap.String("uuid", resume.UUID), zap.Error(err))
		return err
	}

	p.log.Debug("сообщение опубликовано в Kafka", zap.String("uuid", resume.UUID))
	return nil
}

// Close закрывает продюсер, дожидаясь отправки буферизованных сообщений
func (p *Producer) Close() error {
	return p.writer.Close()
}
