package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/app/resize"
	_ "github.com/aaronland/go-image/common"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := resize.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to transform images, %v", err)
	}
}
