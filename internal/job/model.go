package job

import "time"

type Job struct {
	CreatedAt time.Time
	ID        string
	Payload   []byte
	Status    Status
	Retry     int
	MaxRetry  int
}
