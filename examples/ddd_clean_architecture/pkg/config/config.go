package config

// Config holds application configuration
type Config struct {
	DatabaseURL string
	Port        int
	Environment string
}

// NewConfig creates a new configuration
func NewConfig() Config {
	return Config{
		DatabaseURL: "postgres://localhost:5432/myapp",
		Port:        8080,
		Environment: "development",
	}
}

// IsDevelopment returns true if running in development mode
func (c Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c Config) IsProduction() bool {
	return c.Environment == "production"
}
