package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/solrac97gr/goarchtest/test/clean_architecture/application"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	userService *application.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserResponse represents the common response structure for user operations
type UserResponse struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleCreateUser handles the POST /users endpoint
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createReq := application.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
	}

	result, err := h.userService.CreateUser(createReq)
	if err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := UserResponse{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
	}

	respondWithJSON(w, response, http.StatusCreated)
}

// HandleGetUser handles the GET /users/{id} endpoint
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple URL parsing, in a real app, use a router
	id := r.URL.Path[len("/users/"):]
	if id == "" {
		respondWithError(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	result, err := h.userService.GetUserByID(application.GetUserByIDRequest{ID: id})
	if err != nil {
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	response := UserResponse{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
	}

	respondWithJSON(w, response, http.StatusOK)
}

// respondWithJSON writes a JSON response
func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondWithError writes an error response
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	respondWithJSON(w, ErrorResponse{Error: message}, statusCode)
}
