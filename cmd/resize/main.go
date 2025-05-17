package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/resize"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	err := resize.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to resize images, %v", err)
	}
}
