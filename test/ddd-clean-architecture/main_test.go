package main

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestDDDWithCleanArchitecture tests the Domain-Driven Design with Clean Architecture pattern
func TestDDDWithCleanArchitecture(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Create a new Types instance for the project
	types := goarchtest.InPath(projectPath)

	// Define the domains and namespaces for DDD with Clean Architecture
	domains := []string{"user", "order"} // Our bounded contexts
	sharedNamespace := "shared"          // Shared kernel
	pkgNamespace := "pkg"                // Public packages

	// Create the DDD with Clean Architecture pattern
	dddPattern := goarchtest.DDDWithCleanArchitecture(domains, sharedNamespace, pkgNamespace)

	// Validate the pattern against our codebase
	results := dddPattern.Validate(types)

	// Check each validation result
	for i, result := range results {
		t.Run(result.RuleDescription, func(t *testing.T) {
			if !result.IsSuccessful {
				t.Errorf("DDD Clean Architecture Rule %d failed: %s", i+1, result.RuleDescription)
				t.Logf("Failing types:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  - %s in package %s (path: %s)", failingType.Name, failingType.Package, failingType.FullPath)
				}
			} else {
				t.Logf("✅ Rule passed: %s", result.RuleDescription)
			}
		})
	}

	// Additional specific tests for our DDD structure
	t.Run("Domain models should not depend on application layer", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("domain/models").
			ShouldNot().
			HaveDependencyOn("application").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Domain models should not depend on application layer")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
			}
		}
	})

	t.Run("Domain models should not depend on infrastructure", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("domain/models").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Domain models should not depend on infrastructure")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
			}
		}
	})

	t.Run("Domain ports should not depend on infrastructure", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("domain/ports").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Domain ports should not depend on infrastructure")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
			}
		}
	})

	t.Run("Application services should not depend on infrastructure", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("application").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Application services should not depend on infrastructure")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
			}
		}
	})

	t.Run("Infrastructure should implement domain ports", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("infrastructure").
			Should().
			HaveDependencyOn("domain/ports").
			GetResult()

		if !result.IsSuccessful {
			t.Log("ℹ️ Infrastructure components should implement domain ports (this might be expected if no infrastructure is present)")
		}
	})

	t.Run("User domain should not depend on order domain", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("internal/user").
			ShouldNot().
			HaveDependencyOn("internal/order").
			GetResult()

		if !result.IsSuccessful {
			t.Error("User domain should not depend on order domain (bounded context isolation)")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
			}
		}
	})

	t.Run("Order domain should not depend on user domain", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("internal/order").
			ShouldNot().
			HaveDependencyOn("internal/user").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Order domain should not depend on user domain (bounded context isolation)")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s in %s", failingType.Name, failingType.Package)
			}
		}
	})

	t.Run("Both domains can depend on shared kernel", func(t *testing.T) {
		userResult := types.That().
			ResideInNamespace("internal/user").
			Should().
			HaveDependencyOn("shared").
			GetResult()

		orderResult := types.That().
			ResideInNamespace("internal/order").
			Should().
			HaveDependencyOn("shared").
			GetResult()

		if !userResult.IsSuccessful {
			t.Log("ℹ️ User domain doesn't use shared kernel (this might be expected)")
		}

		if !orderResult.IsSuccessful {
			t.Log("ℹ️ Order domain doesn't use shared kernel (this might be expected)")
		}

		// At least one should use shared
		if !userResult.IsSuccessful && !orderResult.IsSuccessful {
			t.Log("ℹ️ Neither domain uses shared kernel - consider if shared types are needed")
		}
	})
}

// TestBoundedContextIsolation specifically tests that bounded contexts are properly isolated
func TestBoundedContextIsolation(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	t.Run("User bounded context isolation", func(t *testing.T) {
		// User context should only depend on shared, not on other bounded contexts
		result := types.That().
			ResideInNamespace("internal/user").
			ShouldNot().
			HaveDependencyOn("internal/order").
			GetResult()

		if !result.IsSuccessful {
			t.Error("User bounded context should not depend on order bounded context")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s depends on order context", failingType.Name)
			}
		}
	})

	t.Run("Order bounded context isolation", func(t *testing.T) {
		// Order context should only depend on shared, not on other bounded contexts
		result := types.That().
			ResideInNamespace("internal/order").
			ShouldNot().
			HaveDependencyOn("internal/user").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Order bounded context should not depend on user bounded context")
			for _, failingType := range result.FailingTypes {
				t.Logf("Violation: %s depends on user context", failingType.Name)
			}
		}
	})
}

// TestCleanArchitectureWithinBoundedContexts tests Clean Architecture rules within each bounded context
func TestCleanArchitectureWithinBoundedContexts(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test Clean Architecture within User bounded context
	t.Run("User context - Clean Architecture", func(t *testing.T) {
		// Domain should not depend on application
		result := types.That().
			ResideInNamespace("internal/user/domain").
			ShouldNot().
			HaveDependencyOn("internal/user/application").
			GetResult()

		if !result.IsSuccessful {
			t.Error("User domain should not depend on user application")
		}

		// Domain should not depend on infrastructure
		result = types.That().
			ResideInNamespace("internal/user/domain").
			ShouldNot().
			HaveDependencyOn("internal/user/infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("User domain should not depend on user infrastructure")
		}

		// Application should not depend on infrastructure
		result = types.That().
			ResideInNamespace("internal/user/application").
			ShouldNot().
			HaveDependencyOn("internal/user/infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("User application should not depend on user infrastructure")
		}
	})

	// Test Clean Architecture within Order bounded context
	t.Run("Order context - Clean Architecture", func(t *testing.T) {
		// Domain should not depend on application
		result := types.That().
			ResideInNamespace("internal/order/domain").
			ShouldNot().
			HaveDependencyOn("internal/order/application").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Order domain should not depend on order application")
		}

		// Domain should not depend on infrastructure
		result = types.That().
			ResideInNamespace("internal/order/domain").
			ShouldNot().
			HaveDependencyOn("internal/order/infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Order domain should not depend on order infrastructure")
		}

		// Application should not depend on infrastructure
		result = types.That().
			ResideInNamespace("internal/order/application").
			ShouldNot().
			HaveDependencyOn("internal/order/infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Error("Order application should not depend on order infrastructure")
		}
	})
}
