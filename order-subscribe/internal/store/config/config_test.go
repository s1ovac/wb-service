package config

import (
	"testing"
)

func TestConfigToStorage(t *testing.T) {
	cf := NewStorageConfig()
	if cf == nil {
		t.Fail()
	}
}
