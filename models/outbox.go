package models

type Outbox struct {
	ID            int64
	EventType     string
	EventData     []byte
	AggregateID   string
	AggregateType string
}
