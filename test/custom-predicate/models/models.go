package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Product represents a product in the system
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Order represents an order in the system
type Order struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Products []Product `json:"products"`
	Total    float64   `json:"total"`
}

// Settings represents application settings
type Settings struct {
	DatabaseURL string `json:"database_url"`
	APIKey      string `json:"api_key"`
}
