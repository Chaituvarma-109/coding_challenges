package main

import (
	"os/exec"
	"testing"
)

func TestNoFile(t *testing.T) {
	file := "text.txt"
	err := exec.Command("./bin/cchead %s", file).Run()

	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestWithFile(t *testing.T) {
	err := exec.Command("./bin/cchead").Run()

	if err != nil {
		t.Errorf("err: %s", err)
	}
}
