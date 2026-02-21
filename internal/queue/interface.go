package queue

import "context"

type Queue interface {
	ConsumeWithContext(ctx context.Context) (string, error)
	PublishWithContext(ctx context.Context) error
}
