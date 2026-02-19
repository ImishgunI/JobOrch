package job

import "time"

type Job struct {
	ID        string
	Payload   []byte
	Status    Status
	Retry     int
	MaxRetry  int
	CreatedAt time.Time
}
