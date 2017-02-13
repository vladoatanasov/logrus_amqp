# AMQP Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

Usage

```go
package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/vladoatanasov/logrus_amqp"
)

func main() {
	log := logrus.New()
  
  	hook := logrus_amqp.NewAMQPHook("127.0.0.1:5672", "guest", "guest", "exchange-rabbitmq", "routing-key")
	// with virtual host
	hook := logrus_amqp.NewAMQPHookWithType("127.0.0.1:5672", "guest", "guest", "exchange-rabbitmq", "virtualHost" "routing-key")
	log.Hooks.Add(hook)
	
}

func doWork() {
  err := some_useful_func()
  
  if err != nil {
		log.WithFields(logrus.Fields{
			"topic": "some_useful_func",
		}).Error(err)
	}
}
```

With this hook, you can easily send logs to the ELK stack, using rabbitmq as a message broker. You can find a working docker-compose project [here](https://github.com/vladoatanasov/docker-elk)
