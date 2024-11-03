package entity

import "context"

type MessageQueue interface {
	PublishMessage(ctx context.Context, key, value []byte) error
}
