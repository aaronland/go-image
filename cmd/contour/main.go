package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/contour"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	err := contour.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to contour images, %v", err)
	}
}
