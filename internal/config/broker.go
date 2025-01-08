package config

import (
	"log"

	"github.com/streadway/amqp"
)

type Broker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

func NewBroker(url, exchange string) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err = conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		err = ch.Close()
		if err != nil {
			return nil, err
		}
		err = conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &Broker{
		connection: conn,
		channel:    ch,
		exchange:   exchange,
	}, nil
}

func (b *Broker) Publish(exchangeName string, routingKey string, qName string, body []byte) error {
	err := setupQueue(b.channel, exchangeName, qName, routingKey)
	if err != nil {
		return err
	}

	return b.channel.Publish(
		b.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (b *Broker) Close() {
	if err := b.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
	}
	if err := b.connection.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

func setupQueue(channel *amqp.Channel, exchangeName string, qName string, keyName string) error {
	_, err := channel.QueueDeclare(
		exchangeName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return channel.QueueBind(
		qName,   //"account.create"
		keyName, //"account.create"
		exchangeName,
		false,
		nil,
	)
}
