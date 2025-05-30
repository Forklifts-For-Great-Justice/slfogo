package slfogolib

import (
	"context"
	"fmt"
	"io"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type WriterOutputProcessor struct {
	w io.Writer
}

func NewWriterOutputProcessor(w io.Writer) WriterOutputProcessor {
	return WriterOutputProcessor{w: w}
}

func (op WriterOutputProcessor) Put(ctx context.Context, lp format.LogParts) error {
	host, err := getKey(lp, "tag")
	if err != nil {
		return err
	}

	msg, err := getKey(lp, "content")
	if err != nil {
		return err
	}
	_, err = op.w.Write([]byte(fmt.Sprintf("%s: %s\n", host, msg)))
	return err
}

func (op WriterOutputProcessor) Close() error {
	return nil
}
