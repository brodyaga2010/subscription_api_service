package main

import (
	"context"
	"eff_mobile/internal"
	"log"
	"time"
)

var (
	Version   = "dev"
	BuildTime = time.Now().Format(time.RFC3339)
)

// @title Subscription Aggregation Service API
// @version 1.0
// @description This is a service for aggregating user online subscriptions.
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()

	log.Printf("Version: %v, BuildTime: %v\n", Version, BuildTime)

	if err := internal.Run(ctx); err != nil {
		log.Fatalf("Error running service: %v", err)
	}

	log.Println("gracefully shutdown!")
}
