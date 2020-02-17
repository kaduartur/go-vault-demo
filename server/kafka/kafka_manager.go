package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/kaduartur/go-vault-demo/server"
)

type Manager struct {
	kafka *kafka.Writer
}

func NewManager(k kafka.WriterConfig) *Manager {
	return &Manager{kafka: kafka.NewWriter(k)}
}

func (m *Manager) SendEvent(p server.PaymentEvent) error {
	log.Printf("Sending event  %+v\n", p)
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&p)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(p.PaymentId),
		Value: buf.Bytes(),
	}
	ctx := context.Background()
	err = m.kafka.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
