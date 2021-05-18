package app

import "context"

type Serve interface {
	Start(ctx context.Context) error
}
