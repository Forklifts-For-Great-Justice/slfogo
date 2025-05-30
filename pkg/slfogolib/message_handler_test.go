package slfogolib

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestHandleMessages(t *testing.T) {
	mop := NewMockOutputProcessor()
	lpc := make(syslog.LogPartsChannel)
	msgVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "Forklift",
			Subsystem: "metrics",
			Name:      "message_count",
			Help:      "its a metric",
		},
		[]string{"server", "service"},
	)

	nmh := NewMessageHandler(mop, lpc, msgVec)
	ctx := context.Background()
	go nmh.HandleMessages(ctx)

	lpc <- format.LogParts{
		"content": "foo",
		"tag":     "test",
	}
	ctx.Done()

	got := mop.GetMessages()[0]
	want := "test: foo"
	if got != want {
		t.Errorf("Got: %s, Want: %s", got, want)
	}
}

func TestHandleMessagesCloseChan(t *testing.T) {
	mop := NewMockOutputProcessor()
	lpc := make(syslog.LogPartsChannel)
	msgVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "Forklift",
			Subsystem: "metrics",
			Name:      "message_count",
			Help:      "its a metric",
		},
		[]string{"server", "service"},
	)

	nmh := NewMessageHandler(mop, lpc, msgVec)
	ctx := context.Background()
	go nmh.HandleMessages(ctx)

	close(lpc)

}

func TestHandleMessagesCloseCtx(t *testing.T) {
	mop := NewMockOutputProcessor()
	lpc := make(syslog.LogPartsChannel)
	msgVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "Forklift",
			Subsystem: "metrics",
			Name:      "message_count",
			Help:      "its a metric",
		},
		[]string{"server", "service"},
	)

	nmh := NewMessageHandler(mop, lpc, msgVec)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	wg := new(sync.WaitGroup)
	go func() {
		wg.Add(1)
		nmh.HandleMessages(ctx)
		wg.Done()
	}()
	time.Sleep(3 * time.Second)
	cancel()
	ctx.Done()
	wg.Wait()

}
