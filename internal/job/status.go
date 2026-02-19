package job

type Status int

const (
	PENDING Status = iota
	PROCESSING
	SUCCEEDED
	FAILED
)
