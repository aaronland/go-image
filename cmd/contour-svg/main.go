package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/contour/svg"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	err := svg.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to contour images, %v", err)
	}
}
