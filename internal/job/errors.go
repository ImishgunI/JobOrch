package job

import "errors"

// Errors for MessageQueue
var (
	ErrNoMessageInQueue = errors.New("queue have no messages")
)
