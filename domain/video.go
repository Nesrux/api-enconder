package domain

import "time"

type Video struct {
	ID         string
	ResourceID string
	filePath   string
	createdAt  time.Time
}
