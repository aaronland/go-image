package transform

import (
	"context"
	"testing"
)

func TestNewTransformation(t *testing.T) {

	ctx := context.Background()

	_, err := NewTransformation(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create new 'null://' transformation, %v", err)
	}

	_, err = NewTransformation(ctx, "fail://")

	if err == nil {
		t.Fatalf("Expected creation of 'fail://' transformation to fail")
	}
}
