package slfogolib

import (
	"context"
	"fmt"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type MockOutputProcessor struct {
	msgs []string
}

func NewMockOutputProcessor() *MockOutputProcessor {
	return &MockOutputProcessor{msgs: make([]string, 0)}
}

func (mop *MockOutputProcessor) Put(ctx context.Context, lp format.LogParts) error {
	host, err := getKey(lp, "tag")
	if err != nil {
		return err
	}

	content, err := getKey(lp, "content")
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("%s: %s", host, content)
	mop.msgs = append(mop.msgs, msg)
	return nil
}

func (mop *MockOutputProcessor) Close() error {
	return nil
}

func (mop *MockOutputProcessor) GetMessages() []string {
	return mop.msgs
}
