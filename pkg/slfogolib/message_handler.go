package slfogolib

import (
	"context"
	"log/slog"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type MessageHandler struct {
	op     OutputProcessor
	lpChan syslog.LogPartsChannel
	mth    *MetricHolder
}

func NewMessageHandler(op OutputProcessor, lpc syslog.LogPartsChannel, mth *MetricHolder) *MessageHandler {
	return &MessageHandler{op: op, lpChan: lpc, mth: mth}
}

func (mh *MessageHandler) updateMetrics(lp format.LogParts) error {
	host, err := getKey(lp, "hostname")
	if err != nil {
		return err
	}
	service, err := getKey(lp, "tag")
	if err != nil {
		return err
	}
	mh.mth.AddMetric([2]string{host, service})

	return nil
}

func (mh *MessageHandler) HandleMessages(ctx context.Context) {
mainLoop:
	for {
		select {
		case lp, ok := <-mh.lpChan:
			if !ok {
				break mainLoop
			} else {
				mh.updateMetrics(lp)
				if err := mh.op.Put(ctx, lp); err != nil {
					slog.ErrorContext(ctx, "Put generated error", "error", err.Error())
				}
			}
		case <-ctx.Done():
			break mainLoop
		}

	}
}
