package slfogolib

import (
	"bytes"
	"context"
	"testing"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestNewWriterOutputProcessor(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	nwop := NewWriterOutputProcessor(buf)
	err := nwop.Put(
		context.TODO(),
		format.LogParts{"content": "test", "tag": "test"},
	)
	if err != nil {
		t.Fatal(err)
	}
}
