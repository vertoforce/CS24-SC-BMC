package bmc

import (
	"context"
	"os"
	"testing"
)

func TestSetPowerState(t *testing.T) {
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
