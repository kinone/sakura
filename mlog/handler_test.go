package mlog

import (
	"testing"
)

func TestLevelFilter(t *testing.T) {
	f := LevelFilter("error+")

	if f(&Record{level: Info}) {
		t.Error("failed")
	}

	if f(&Record{level:Warning}) {
		t.Error("failed")
	}

	if !f(&Record{level: Error}) {
		t.Error("failed")
	}

	if !f(&Record{level: Critical}) {
		t.Error("failed")
	}

	if !f(&Record{level: Emergency}) {
		t.Error("failed")
	}
}
