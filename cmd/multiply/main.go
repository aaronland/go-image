package main

import (
	"context"
	"log"

	"github.com/aaronland/go-image/v2/app/multiply"
	_ "github.com/aaronland/go-image/v2/common"
)

func main() {

	ctx := context.Background()
	err := multiply.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to multiply images, %v", err)
	}
}
