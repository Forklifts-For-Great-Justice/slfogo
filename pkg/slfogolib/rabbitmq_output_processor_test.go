package slfogolib

import (
	"bytes"
	"context"
	"testing"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestNewRabbitMQOutputProcessor(t *testing.T) {
	wantUser := "guest"
	wantPass := "guest"
	wantHostname := "rabbit.example.com"
	wantPort := 5672
	wantExchangeName := "testExchange"

	got := NewRabbitMQOutputProcessor(
		wantUser,
		wantPass,
		wantHostname,
		wantPort,
		wantExchangeName,
	)

	if got == nil {
		t.Fatal("Unexpected nil value")
	}

	if got.username != wantUser {
		t.Errorf("Got: %s, Want: %s", got.username, wantUser)
	}

	if got.password != wantPass {
		t.Errorf("Got: %s, Want: %s", got.password, wantPass)
	}

	if got.hostname != wantHostname {
		t.Errorf("Got: %s, Want: %s", got.hostname, wantHostname)
	}

	if got.port != wantPort {
		t.Errorf("Got: %d, Want: %d", got.port, wantPort)
	}

	if got.exchangeName != wantExchangeName {
		t.Errorf("Got: %s, Want: %s", got.exchangeName, wantExchangeName)
	}

	if got.Connect == nil {
		t.Errorf("Got nil, Want: not nil")
	}
}

func TestRabbitBuildConnectString(t *testing.T) {
	rab := NewRabbitMQOutputProcessor(
		"guest",
		"guest",
		"rabbit.example.com",
		5672,
		"testExchange",
	)

	got := rab.buildConnectString()
	want := "amqp://guest:guest@rabbit.example.com:5672"
	if got != want {
		t.Errorf("BuildConnectString() Got: %s, Want: %s", got, want)
	}
}

func TestBuildRoutingKey(t *testing.T) {
	tests := []struct {
		name    string
		lp      format.LogParts
		want    string
		wantErr bool
	}{
		{
			name: "Get Routing Key",
			lp: format.LogParts{
				"tag": "test-service",
			},
			want:    "syslog.test-service",
			wantErr: false,
		},
		{
			name:    "Get error",
			lp:      format.LogParts{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := buildRoutingKey(tc.lp)
			if (gotErr != nil) != tc.wantErr {
				t.Errorf("buildRoutingKey(%v) got unexpected error %v", tc.lp, gotErr)
			}

			if got != tc.want {
				t.Errorf("buildRoutingKey(%v): got: %s, want: %s", tc.lp, got, tc.want)
			}
		})
	}

}

func TestPublishMessage(t *testing.T) {
	rab := NewRabbitMQOutputProcessor(
		"guest",
		"guest",
		"localhost",
		5672,
		"testExchange",
	)
	rab.Connect = rab.mockConnect
	err := rab.Connect()
	if err != nil {
		t.Errorf("rab.Connect(): Unexpected error %v", err)
	}

	lp := format.LogParts{"tag": "test-service", "content": "test string"}
	if err := rab.Put(context.TODO(), lp); err != nil {
		t.Errorf("rab.Put(), unexpected error: %v", err)
	}

	got := rab.ch.(*MockPublisher).entries[0]

	if got.exchangeName != "testExchange" {
		t.Errorf("got.exchangeName: %s want: testExchange", got.exchangeName)
	}

	if got.routingKey != "syslog.test-service" {
		t.Errorf("got.routingKey: %s want: syslog.test-service", got.routingKey)
	}

	if !bytes.Equal(got.body, []byte("test string")) {
		t.Errorf("got.body: %s want: []byte(\"test string\")", got.body)
	}

}
