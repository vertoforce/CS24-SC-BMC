package bmc

import (
	"context"
	"os"
	"testing"
)

func TestTemp(t *testing.T) {
	c, err := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Error(err)
		return
	}

	result, err := c.GetTemperature(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if len(result) == 0 {
		t.Errorf("No temperatures found")
	}
}
