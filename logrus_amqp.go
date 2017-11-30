package logrus_amqp

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AMQPHook struct {
	AMQPServer   string
	Username     string
	Password     string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	VirtualHost  string
	Mandatory    bool
	Immediate    bool
	Durable      bool
	Internal     bool
	NoWait       bool
	AutoDeleted  bool
}

func NewAMQPHook(server, username, password, exchange, routingKey string) *AMQPHook {
	return NewAMQPHookWithType(server, username, password, exchange, "direct", "", routingKey)
}

func NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey string) *AMQPHook {
	hook := AMQPHook{}

	hook.AMQPServer = server
	hook.Username = username
	hook.Password = password
	hook.Exchange = exchange
	hook.ExchangeType = exchangeType
	hook.Durable = true
	hook.RoutingKey = routingKey
	hook.VirtualHost = virtualHost

	return &hook
}

// Fire is called when an event should be sent to the message broker
func (hook *AMQPHook) Fire(entry *logrus.Entry) error {
	dialURL := fmt.Sprintf("amqp://%s:%s@%s/%s", hook.Username, hook.Password, hook.AMQPServer, hook.VirtualHost)
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

	body, err := entry.String()
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
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}

	return nil
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
