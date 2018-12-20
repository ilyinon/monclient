package main

import (
	"testing"
)

func TestGetService(t *testing.T) {
	data := []struct {
		hostname string
		want     string
	}{
		{"qafrmapp611.scl", "FDM"},
		{"qacore111.scl", "VCC"},
		{"qasccnas612a.scl", "SCC"},
		{"nocore111.scl", "fail"},
	}

	for _, td := range data {
		res := get_service(td.hostname)
		if res != td.want {
			t.Errorf("passed %s, res %s, but want %s", td.hostname, res, td.want)
		}
	}
}

