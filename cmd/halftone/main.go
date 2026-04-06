package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/halftone"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	err := halftone.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to halftone images, %v", err)
	}
}
