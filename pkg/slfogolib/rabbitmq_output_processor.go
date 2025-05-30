package slfogolib

import (
	"context"
	"fmt"
	"io"

	"github.com/prometheus/client_golang/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type Publisher interface {
	PublishWithContext(
		context.Context,
		string,
		string,
		bool,
		bool,
		amqp.Publishing,
	) error
	io.Closer
}

func (rab *RabbitMQOutputProcessor) mockConnect() error {
	ch := new(MockPublisher)
	ch.entries = []entry{}
	rab.ch = ch
	rab.conn = new(mockConn)

	return nil
}

type RabbitMQOutputProcessor struct {
	username     string
	password     string
	hostname     string
	port         int
	exchangeName string
	conn         io.Closer
	ch           Publisher
	Connect      func() error
	msgGauge     prometheus.GaugeVec
}

func NewRabbitMQOutputProcessor(
	username string,
	password string,
	hostname string,
	port int,
	exchangeName string,
) *RabbitMQOutputProcessor {

	rab := &RabbitMQOutputProcessor{
		username:     username,
		password:     password,
		hostname:     hostname,
		port:         port,
		exchangeName: exchangeName,
	}

	rab.Connect = rab.connect

	return rab
}

func (rab *RabbitMQOutputProcessor) buildConnectString() string {
	if rab.username == "" {
		return fmt.Sprintf(
			"amqp://%s:%d/%%2f",
			rab.hostname,
			rab.port,
		)
	}

	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		rab.username,
		rab.password,
		rab.hostname,
		rab.port,
	)
}

func buildRoutingKey(lp format.LogParts) (string, error) {
	host, err := getKey(lp, "tag")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("syslog.%s", host), nil
}

func (rab *RabbitMQOutputProcessor) connect() error {
	conn, err := amqp.Dial(rab.buildConnectString())
	if err != nil {
		return err
	}
	rab.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	rab.ch = ch

	if err := ch.ExchangeDeclare(
		rab.exchangeName,
		"topic", // exchange type
		true,    // durable
		false,   // autoDelete
		false,   // internal
		false,   // nowait
		nil,     // args
	); err != nil {
		return err
	}

	return nil
}

func (rab *RabbitMQOutputProcessor) Put(ctx context.Context, lp format.LogParts) error {
	routingKey, err := buildRoutingKey(lp)
	if err != nil {
		return err
	}

	/*
		host, err := getKey(lp, "tag")
		if err != nil {
			return err
		}
	*/

	content, err := getKey(lp, "content")
	if err != nil {
		return err
	}

	err = rab.ch.PublishWithContext(
		ctx,
		rab.exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Body: []byte(content),
		},
	)
	return err
}

func (rab *RabbitMQOutputProcessor) Close() error {
	if err := rab.ch.Close(); err != nil {
		return err
	}

	if err := rab.conn.Close(); err != nil {
		return err
	}

	return nil
}
