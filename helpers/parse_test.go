package helpers_test

import (
	"testing"

	"github.com/edulustosa/go-pay/helpers"
)

func TestParseDocument(t *testing.T) {
	testCases := []struct {
		document string
		want     bool
	}{
		{"529.982.247-25", true},
		{"52998224725", true},
		{"168.995.350-09", true},
		{"16899535009", true},
		{"529.982.247-26", false},
		{"70696857189", true},
		{"11111111111", false},
	}

	for _, tc := range testCases {
		err := helpers.ParseDocument(tc.document)
		if err != nil && tc.want {
			t.Errorf("ParseDocument(%s) got %v, want nil", tc.document, err)
		}
		if err == nil && !tc.want {
			t.Errorf("ParseDocument(%s) got nil, want error", tc.document)
		}
	}
}
