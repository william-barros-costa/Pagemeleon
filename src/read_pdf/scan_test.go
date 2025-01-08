package main

import (
	"testing"
)

func TestCreateScanner(t *testing.T) {
	t.Run("Test create Scan", func(t *testing.T) {
		scan := scanState{}
		if scan.isHex != false {
			t.Errorf("Expected isHex to be false, got %v", scan.isHex)
		}
	})
	t.Run("Other test", func(t *testing.T) {
		scan := scanState{}
		if scan.isValue != false {
			t.Errorf("Expected isValue to be false, got %v", scan.isValue)
		}
	})
}
