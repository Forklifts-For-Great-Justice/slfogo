package slfogolib

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type mockConn struct {
	status bool
}

func (mc *mockConn) Close() error {
	mc.status = false

	return nil
}

type entry struct {
	routingKey   string
	exchangeName string
	body         []byte
}

type MockPublisher struct {
	entries []entry
}

func (mp *MockPublisher) PublishWithContext(
	_ context.Context,
	exchangeName string,
	routingKey string,
	_ bool,
	_ bool,
	pub amqp.Publishing,
) error {
	mp.entries = append(mp.entries, entry{exchangeName: exchangeName, routingKey: routingKey, body: pub.Body})

	return nil
}

func (mp *MockPublisher) Close() error {
	return nil
}
