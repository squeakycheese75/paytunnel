package repository

import "time"

type Event struct {
	EventID    int64
	DeliveryID string
	EventName  string
	TargetUrl  string
	Secret     string
	BodyJson   string
	CreatedAt  time.Time
}
