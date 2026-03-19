package entities

type Event struct {
	EventID    int64
	DeliveryID string
	EventName  string
	TargetUrl  string
	BodyJson   string
	CreatedAt  string
}
