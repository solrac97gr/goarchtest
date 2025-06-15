package shared

import (
	"crypto/rand"
	"fmt"
	"time"
)

// ID represents a shared identifier type across all domains
type ID struct {
	Value string
}

// NewID creates a new unique identifier
func NewID() ID {
	b := make([]byte, 16)
	rand.Read(b)
	return ID{Value: fmt.Sprintf("%x", b)}
}

// Event represents a domain event
type Event interface {
	GetEventType() string
	GetTimestamp() time.Time
}

// AggregateRoot represents a domain aggregate root
type AggregateRoot interface {
	GetID() ID
	GetEvents() []Event
	ClearEvents()
}
