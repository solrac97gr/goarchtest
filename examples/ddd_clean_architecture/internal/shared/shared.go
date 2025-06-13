package shared

// EntityStatus represents the status of any entity in the system
// This is part of the shared kernel - common concepts used across domains
type EntityStatus int

const (
	StatusActive EntityStatus = iota
	StatusInactive
	StatusPending
	StatusArchived
)

// String returns the string representation of the status
func (s EntityStatus) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	case StatusPending:
		return "pending"
	case StatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// AuditInfo provides common audit information for entities
type AuditInfo struct {
	CreatedAt string
	UpdatedAt string
	CreatedBy string
	UpdatedBy string
}

// DomainEvent represents a domain event that can occur in any domain
type DomainEvent interface {
	EventName() string
	OccurredAt() string
	AggregateID() string
}
