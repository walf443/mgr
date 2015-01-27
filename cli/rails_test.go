package main

import (
	"os"
	"testing"
)

func TestDetectFramework(t *testing.T) {
	var f RubyOnRailsFramework
	if f.DetectFramework() != false {
		t.Errorf("this directory should not be detected as rails framework")
	}
	err := os.Chdir("railssample")
	if err != nil {
		t.Errorf("Can't change dir to railssample")
	}

	if !f.DetectFramework() {
		t.Errorf("this directory should be detected as rails framework")
	}
}
