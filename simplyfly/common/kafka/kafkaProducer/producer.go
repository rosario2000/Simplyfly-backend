package kafkaProducer

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"simplifly/common/constants"
)

var Producer *kafka.Writer

func Init() {
	mechanism, _ := scram.Mechanism(scram.SHA256, "b3V0Z29pbmctc3RpbmdyYXktNzk4NyQ0E8IMaeo_CZgv-TK_uRmNDDryWeQljk4", "ZjdhNzA5NTgtOTRmNS00ZjFlLWFmODMtYTk1YWUwZDE1OGIz")
	Producer = &kafka.Writer{
		Topic:                  "user_updates_topic",
		Addr:                   kafka.TCP(constants.KafkaBroker),
		AllowAutoTopicCreation: true,
		Async:                  false,
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}
}

func SendUpdateToUser(ctx context.Context, payload interface{}) (err error) {
	if Producer == nil {
		Init()
	}
	var byteRes []byte
	byteRes, err = json.Marshal(payload)
	if err != nil {
		return
	}
	log.Info("SendEventMsgToKafka : ", string(byteRes))
	return Producer.WriteMessages(ctx, kafka.Message{Value: byteRes})
}
