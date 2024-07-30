package kafkaConsumer

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"os"
	"os/signal"
	"syscall"
)

var broadcast = make(chan []byte)

func newConsumer(topic string) *kafka.Reader {
	mechanism, _ := scram.Mechanism(scram.SHA512, "b3V0Z29pbmctc3RpbmdyYXktNzk4NyQ0E8IMaeo_CZgv-TK_uRmNDDryWeQljk4", "ZjdhNzA5NTgtOTRmNS00ZjFlLWFmODMtYTk1YWUwZDE1OGIz")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"outgoing-stingray-7987-eu2-kafka.upstash.io:9092"},
		Topic:   topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
		GroupID:  "simplyfly-consumer-group",
		MaxBytes: 10e6, // 10MB
	})
}

func StartConsumer(ctx context.Context, topic string, kafkaMessageProcessor func([]byte)) (err error) {
	reader := newConsumer(topic)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			log.Infof("Caught signal %v: terminating", sig)
			run = false
		default:
			message, rErr := reader.ReadMessage(ctx)
			if rErr != nil && rErr.Error() != "fetching message: unexpected EOF" {
				log.Errorj(log.JSON{
					"topic":   message.Topic,
					"message": fmt.Sprintf("read error"),
					"error":   rErr.Error(),
				})
				break
			}
			kafkaMessageProcessor(message.Value)
		}
	}
	log.Infoj(log.JSON{"message": "Close kafka consumer", "topic": topic})
	return reader.Close()
}
