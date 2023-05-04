package main

import (
	"GoScan"
	"testing"
)

func TestScanPort(t *testing.T) {
	got := GoScan.ScanPort("localhost", 22)
	var want GoScan.ScanResult
	want.Port = 22
	want.State = true

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
