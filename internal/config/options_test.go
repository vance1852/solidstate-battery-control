package config

import "testing"

func TestOptions(t *testing.T) {
	o := DefaultOptions()
	if o.Validate() != nil {
		t.Fatal("defaults")
	}
	o.MaxBodyBytes = 1
	if o.Validate() == nil {
		t.Fatal("small body")
	}
	o = LoadOptions()
	if o.MaxBodyBytes < 1024 {
		t.Fatal("load")
	}
}
