package GoScan

import (
	"GoScan"
	"testing"
)

func TestScanPort(t *testing.T) {
	got := ScanHost("localhost", 22)
	var want ScanResult
	want.Port = 22
	want.State = true

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
