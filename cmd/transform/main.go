package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/transform"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := transform.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to transform images, %v", err)
	}
}
