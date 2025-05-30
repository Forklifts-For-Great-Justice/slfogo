package main

import (
	"log/slog"
	"math/rand/v2"
	"time"
)

func main() {
	slog.Info("testbin starting")
	for {
		sleepTime := rand.Int64N(50)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		slog.Info("ping")
	}
}
