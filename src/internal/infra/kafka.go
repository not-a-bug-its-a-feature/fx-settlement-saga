package infra

import (
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// Brokerの文字列をパースし指定のtopicにメッセージを発行用のWriter
func NewKafkaWriter(brokersCSV string, topic string) (*kafka.Writer, error) {
	if strings.TrimSpace(brokersCSV) == "" {
		return nil, fmt.Errorf("kafka brokers is required")
	}

	if strings.TrimSpace(topic) == "" {
		return nil, fmt.Errorf("kafka topic is required")
	}

	brokers := strings.Split(brokersCSV, ",")

	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	writer := &kafka.Writer{
		Addr: kafka.TCP(brokers...),
		Topic: topic,
		Balancer: &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return writer, nil
}