package bmc

import (
	"context"
	"os"
	"testing"
)

// TestSetPowerState physically power-cycles the server; it only runs when the
// BMC hardware env vars are set.
func TestSetPowerState(t *testing.T) {
	if os.Getenv("IP") == "" {
		t.Skip("BMC hardware env not set")
	}
	c, err := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Error(err)
		return
	}

	err = c.Start(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	err = c.SetPowerState(context.Background(), PowerStop)
	if err != nil {
		t.Error(err)
	}
}
