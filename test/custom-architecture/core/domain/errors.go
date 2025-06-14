package domain

// DomainError represents a domain-specific error
type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// BusinessRule represents a business rule validation
type BusinessRule interface {
	Validate() error
	Description() string
}
