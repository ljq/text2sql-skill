package utils

import (
	"context"
	"time"
)

func WithResourceLimits(ctx context.Context, memoryMB int, timeoutMs int) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)

	// Memory limit enforcement would be implemented at the OS level in production
	// This is a placeholder for the concept

	return ctx, func() {
		cancel()
		// Cleanup memory resources if needed
	}
}
