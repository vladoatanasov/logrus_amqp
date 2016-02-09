package logrus_amqp

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AMQPHook struct {
	AMQPServer   string
	Username     string
	Password     string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Mandatory    bool
	Immediate    bool
	Durable      bool
	Internal     bool
	NoWait       bool
	AutoDeleted  bool
}

func NewAMQPHook(server, username, password, exchange, routingKey string) *AMQPHook {
	hook := AMQPHook{}

	hook.AMQPServer = server
	hook.Username = username
	hook.Password = password
	hook.Exchange = exchange
	hook.ExchangeType = "direct"
	hook.Durable = true
	hook.RoutingKey = routingKey

	return &hook
}

// Fire is called when an event should be sent to the message broker
func (hook *AMQPHook) Fire(entry *logrus.Entry) error {
	dialURL := fmt.Sprintf("amqp://%s:%s@%s/", hook.Username, hook.Password, hook.AMQPServer)
	conn, err := amqp.Dial(dialURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		hook.Exchange,
		hook.ExchangeType,
		hook.Durable,
		hook.AutoDeleted,
		hook.Internal,
		hook.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := getBytes(entry.Data)
	if err != nil {
		return err
	}

	err = ch.Publish(
		hook.Exchange,
		hook.RoutingKey,
		hook.Mandatory,
		hook.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Levels is available logging levels.
func (hook *AMQPHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
