package bmc

import (
	"context"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	_, err := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		t.Error(err)
	}
}
