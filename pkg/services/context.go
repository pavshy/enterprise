package services

import (
	"context"
	"time"
)

const defaultOperationTimeout = time.Second * 20

func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultOperationTimeout)
}
