package tests

import (
	. "bot/models"
	"testing"
)

func TestWarid(t *testing.T) {
	id, err := GetNextWarId()
	if err != nil {
		t.Fatalf("warid err:%s\n", err.Error())
	}
	t.Logf("warid is %d\n", id)
}
