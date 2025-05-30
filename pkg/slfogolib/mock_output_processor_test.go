package slfogolib

import (
	"context"
	"slices"
	"testing"

	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestNewMockOutputProcessor(t *testing.T) {
	mop := NewMockOutputProcessor()
	if mop == nil {
		t.Fatal("Got nil, Want not nil")
	}

	mop.Put(
		context.TODO(),
		format.LogParts{
			"content": "foo",
			"tag":     "test",
		})
	mop.Put(
		context.TODO(),
		format.LogParts{
			"content": "bar",
			"tag":     "test",
		})

	got := mop.GetMessages()
	want := []string{"test: foo", "test: bar"}
	if !slices.Equal(got, want) {
		t.Errorf("Got %v, Want %v", got, want)
	}

}

func TestNewMessageHandler(t *testing.T) {
	mop := NewMockOutputProcessor()
	lpc := make(syslog.LogPartsChannel)

	nmh := NewMessageHandler(mop, lpc, nil)
	if nmh == nil {
		t.Fatal("NewMessageHandler returned nil")
	}

	if nmh.op == nil {
		t.Fatal("nmh.op is nil")
	}

	if nmh.lpChan == nil {
		t.Fatal("nmh.lpChan is nil")
	}

	if nmh.op != mop {
		t.Errorf("Got: nmh.op %v, Want: %v", nmh.op, mop)
	}

	if nmh.lpChan != lpc {
		t.Errorf("Got: nmh.lpChan %v, Want: %v", nmh.lpChan, lpc)
	}
}
