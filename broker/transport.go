package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var AMQP_PROTOCOL = "amqp"

type BrokerConnectionConfig struct {
	Protocol string
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
}

func (c *BrokerConnectionConfig) Dsn() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d%s",
		c.Protocol, c.User, c.Password, c.Host, c.Port, c.Vhost,
	)
}

func (c *BrokerConnectionConfig) OpenTransport() (*amqp.Connection, error) {
	return amqp.Dial(c.Dsn())
}

func DeclareQueue(queueName string, channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}
