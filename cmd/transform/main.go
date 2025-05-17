package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/transform"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	err := transform.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to transform images, %v", err)
	}
}
