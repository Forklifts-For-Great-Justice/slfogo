package slfogolib

import (
	"testing"
	/*
			"context"
			"gopkg.in/mcuadros/go-syslog.v2/format"
		"gopkg.in/mcuadros/go-syslog.v2"
	*/)

func TestGetPort(t *testing.T) {
	tests := []struct {
		name    string
		portVal string
		want    int64
	}{
		{
			name:    "Get actual port value",
			portVal: "8989",
			want:    8989,
		},
		{
			name:    "Get Default port",
			portVal: "",
			want:    9999,
		},
		{
			name:    "Not a string with an int",
			portVal: "foobar9",
			want:    9999,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := GetPort(tc.portVal)

			if got != tc.want {
				t.Errorf("GetPort(%s): Got: %d, Want: %d", tc.portVal, got, tc.want)
			}
		})
	}
}

func TestBuildConnectString(t *testing.T) {
	tests := []struct {
		name    string
		portVal string
		want    string
	}{
		{
			name:    "Get actual port value",
			portVal: "8989",
			want:    "0.0.0.0:8989",
		},
		{
			name:    "Get Default port",
			portVal: "",
			want:    "0.0.0.0:9999",
		},
		{
			name:    "Not a string with an int",
			portVal: "foobar9",
			want:    "0.0.0.0:9999",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildConnectString(tc.portVal)

			if got != tc.want {
				t.Errorf("BuildConnectString(%s): Got: %s, Want: %s", tc.portVal, got, tc.want)
			}
		})
	}
}

func TestBuildServer(t *testing.T) {

	gotSvr, lpc := BuildServer()

	if gotSvr == nil {
		t.Fatal("got nil from gotSvr")
	}

	if lpc == nil {
		t.Fatal("got nil for lpc")
	}

}
