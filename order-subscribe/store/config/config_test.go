package config

import (
	"testing"
)

func TestConfigToStorage(t *testing.T) {
	cf := NewConfig()
	if cf == nil {
		t.Fail()
	}
}
