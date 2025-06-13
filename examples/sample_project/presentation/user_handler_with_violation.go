package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/solrac97gr/goarchtest/examples/sample_project/application"
	// ⚠️  INTENTIONAL CLEAN ARCHITECTURE VIOLATION ⚠️
	// This import creates a Clean Architecture violation!
	// Presentation layer should NOT directly depend on infrastructure layer
	// This is included to demonstrate how GoArchTest detects violations
	"github.com/solrac97gr/goarchtest/examples/sample_project/infrastructure"
)

// UserHandlerWithViolation demonstrates a Clean Architecture violation
// by having the presentation layer directly depend on the infrastructure layer
type UserHandlerWithViolation struct {
	userService *application.UserService
	// This field creates a direct dependency from presentation to infrastructure
	// which violates Clean Architecture - presentation should only depend on application
	userRepo *infrastructure.InMemoryUserRepository
}

// NewUserHandlerWithViolation creates a new UserHandlerWithViolation
// This violates Clean Architecture by creating a direct dependency on infrastructure
func NewUserHandlerWithViolation(userService *application.UserService, userRepo *infrastructure.InMemoryUserRepository) *UserHandlerWithViolation {
	return &UserHandlerWithViolation{
		userService: userService,
		userRepo:    userRepo, // Direct dependency on infrastructure - violation!
	}
}

// GetUserDirectly bypasses the application layer and accesses infrastructure directly
// This is a violation of Clean Architecture
func (h *UserHandlerWithViolation) GetUserDirectly(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	// Violation: Presentation layer directly accessing infrastructure
	// Should go through application layer instead
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
