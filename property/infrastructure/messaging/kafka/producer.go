package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) PublishMessage(ctx context.Context, key, value []byte) error {
	// Write the message to Kafka
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})

	if err != nil {
		log.Printf("failed to publish message: %v", err)
	}
	// Check if the message was successfully written
	return err
}

func (p *KafkaProducer) Close() {
	p.writer.Close()
}

