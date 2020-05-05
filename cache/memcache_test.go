package cache

import (
	"testing"
	"time"
)

func TestMemcache(t *testing.T) {
	mem := NewMemcache("127.0.0.1:11211")
	var err error
	timeoutDuration := 10 * time.Second
	if err = mem.Set("username", "antsbean", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if !mem.IsExist("username") {
		t.Error("IsExist Error")
	}

	name := mem.Get("username").(string)
	if name != "antsbean" {
		t.Error("get Error")
	}

	if err = mem.Delete("username"); err != nil {
		t.Errorf("delete Error , err=%v", err)
	}
}
