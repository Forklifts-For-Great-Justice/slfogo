package slfogolib

import (
	"testing"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestGetKey(t *testing.T) {
	tests := []struct {
		name      string
		lp        format.LogParts
		searchKey string
		want      string
		wantErr   bool
	}{
		{
			name:      "Get Key",
			lp:        format.LogParts{"test": "foo"},
			searchKey: "test",
			want:      "foo",
			wantErr:   false,
		},
		{
			name:      "No key",
			lp:        format.LogParts{},
			searchKey: "test",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "Invalid value",
			lp:        format.LogParts{"test": struct{ a int }{a: 5}},
			searchKey: "test",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := getKey(tc.lp, tc.searchKey)
			if (gotErr != nil) != tc.wantErr {
				t.Errorf("Unexpected error: Get(%v, %s): %v", tc.lp, tc.searchKey, gotErr)
			}

			if got != tc.want {
				t.Errorf("Get(%v, %s) Got: %s, Want: %s", tc.lp, tc.searchKey, got, tc.want)
			}
		})
	}
}
