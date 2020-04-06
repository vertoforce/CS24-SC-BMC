package bmc

import (
	"context"
	"os"
)

func Example() {
	c, _ := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	c.Start(context.Background())
}
