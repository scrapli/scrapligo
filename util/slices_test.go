package util

import (
	"testing"
)

func TestStringSliceContains(t *testing.T) {
	tests := []struct {
		desc string
		ss   []string
		s    string
		want bool
	}{{
		desc: "string-slice-contains-simple-true",
		ss:   []string{"taco", "cat", "race", "car"},
		s:    "taco",
		want: true,
	}, {
		desc: "string-slice-contains-simple-false",
		ss:   []string{"taco", "cat", "race", "car"},
		s:    "potato",
		want: false,
	}}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got := StringSliceContains(tt.ss, tt.s); got != tt.want {
				t.Fatalf("StringSliceContains(%v, %q) failed: got %v, want %v", tt.ss, tt.s, got, tt.want)
			}
		})
	}
}
