package bmc

import (
	"context"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	if os.Getenv("IP") == "" {
		t.Skip("BMC hardware env not set")
	}
	_, err := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Error(err)
	}
}
