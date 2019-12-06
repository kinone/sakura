package mlog

import "testing"

func TestNewLevel(t *testing.T) {
	if NewLevel("info+") != Info|Notice|Warning|Error|Alert|Critical|Emergency {
		t.Error("failed")
	}

	if NewLevel("notice+") != Notice|Warning|Error|Alert|Critical|Emergency {
		t.Error("failed")
	}

	if NewLevel("warning+") != Warning|Error|Alert|Critical|Emergency {
		t.Error("failed")
	}

	if NewLevel("error+") != LevelAll & ^Debug & ^Info & ^Notice & ^Warning {
		t.Error("failed")
	}
}

func TestLevel_String(t *testing.T) {
	if Debug.String() != "[DEBUG]" {
		t.Error("failed")
	}

	if Info.String() != "[INFO]" {
		t.Error("failed")
	}

	if Error.String() != "[ERROR]" {
		t.Error("failed")
	}
}
