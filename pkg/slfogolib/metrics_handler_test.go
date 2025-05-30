package slfogolib

import (
	"testing"
)

func TestNewMetricHolder(t *testing.T) {
	mh := NewMetricHolder()
	if mh == nil {
		t.Fatal("mh is nil, expected not nil")
	}

	if mh.msgVec == nil {
		t.Error("mh.msgVec is nil, expected not nil")
	}
}
