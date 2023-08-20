package localcache

import (
	"testing"
	"time"
)

func TestSetKV(t *testing.T) {
	SetKV("key", "value", "value111", time.Minute)
	for i := 0; i < 100; i++ {
		t.Log("get", GetStringValue("key", "value"))
	}
}
