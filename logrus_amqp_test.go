package logrus_amqp

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConnection struct {
	mock.Mock
}

func TestNewAMQPHook(t *testing.T) {
	hook := NewAMQPHook("server.name", "username", "pass", "exchange", "route")

	if assert.NotNil(t, hook) {
		assert.Equal(t, "server.name", hook.AMQPServer)
		assert.Equal(t, "username", hook.Username)
		assert.Equal(t, "pass", hook.Password)
		assert.Equal(t, "exchange", hook.Exchange)
		assert.Equal(t, "direct", hook.ExchangeType)
		assert.Equal(t, "", hook.VirtualHost)
		assert.Equal(t, "route", hook.RoutingKey)
	}
}

func TestNewAMQPHookWithType(t *testing.T) {
	hook := NewAMQPHookWithType("server.name", "username", "pass", "exchange", "fanout", "virtual", "route")

	if assert.NotNil(t, hook) {
		assert.Equal(t, "server.name", hook.AMQPServer)
		assert.Equal(t, "username", hook.Username)
		assert.Equal(t, "pass", hook.Password)
		assert.Equal(t, "exchange", hook.Exchange)
		assert.Equal(t, "fanout", hook.ExchangeType)
		assert.Equal(t, "virtual", hook.VirtualHost)
		assert.Equal(t, "route", hook.RoutingKey)
	}
}

func TestLevels(t *testing.T) {
	hook := NewAMQPHook("server.name", "username", "pass", "exchange", "route")
	levels := hook.Levels()

	if assert.NotNil(t, levels) {
		assert.Len(t, levels, 7)
		assert.Contains(t, levels, logrus.PanicLevel)
		assert.Contains(t, levels, logrus.FatalLevel)
		assert.Contains(t, levels, logrus.ErrorLevel)
		assert.Contains(t, levels, logrus.WarnLevel)
		assert.Contains(t, levels, logrus.InfoLevel)
		assert.Contains(t, levels, logrus.DebugLevel)
		assert.Contains(t, levels, logrus.TraceLevel)
	}
}
