package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestCustomPredicates(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test 1: Find all handlers/controllers (demonstrating HaveNameEndingWith)
	t.Run("Find all handlers and controllers", func(t *testing.T) {
		handlers := types.That().
			HaveNameEndingWith("Handler").
			GetResult()

		controllers := types.That().
			HaveNameEndingWith("Controller").
			GetResult()

		// For positive matches, if successful, we found types
		if handlers.IsSuccessful {
			fmt.Printf("Found %d handlers:\n", len(handlers.FailingTypes))
			for _, handler := range handlers.FailingTypes {
				fmt.Printf("  - %s in %s\n", handler.Name, handler.Package)
			}
		} else {
			fmt.Println("No handlers found")
		}

		if controllers.IsSuccessful {
			fmt.Printf("Found %d controllers:\n", len(controllers.FailingTypes))
			for _, controller := range controllers.FailingTypes {
				fmt.Printf("  - %s in %s\n", controller.Name, controller.Package)
			}
		} else {
			fmt.Println("No controllers found")
		}
	})

	// Test 2: Find all services (demonstrating HaveNameEndingWith)
	t.Run("Find all services", func(t *testing.T) {
		result := types.That().
			HaveNameEndingWith("Service").
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d services:\n", len(result.FailingTypes))
			for _, service := range result.FailingTypes {
				fmt.Printf("  - %s in %s\n", service.Name, service.Package)
			}
		}
	})

	// Test 3: Find types that start with specific prefix (demonstrating HaveNameStartingWith)
	t.Run("Find types starting with User", func(t *testing.T) {
		result := types.That().
			HaveNameStartingWith("User").
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d types starting with 'User':\n", len(result.FailingTypes))
			for _, userType := range result.FailingTypes {
				fmt.Printf("  - %s in %s\n", userType.Name, userType.Package)
			}
		}
	})

	// Test 4: Find types matching regex pattern (demonstrating HaveNameMatching)
	t.Run("Find types matching regex pattern", func(t *testing.T) {
		result := types.That().
			HaveNameMatching(".*Helper").
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d helper types:\n", len(result.FailingTypes))
			for _, helper := range result.FailingTypes {
				fmt.Printf("  - %s in %s\n", helper.Name, helper.Package)
			}
		}
	})

	// Test 5: Find all struct types (demonstrating BeStruct)
	t.Run("Find all struct types", func(t *testing.T) {
		result := types.That().
			BeStruct().
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d struct types:\n", len(result.FailingTypes))
			for _, structType := range result.FailingTypes {
				fmt.Printf("  - %s in %s\n", structType.Name, structType.Package)
			}
		}
	})

	// Test 6: Combine multiple predicates with And
	t.Run("Find services that are structs", func(t *testing.T) {
		result := types.That().
			HaveNameEndingWith("Service").
			And().
			BeStruct().
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d service structs:\n", len(result.FailingTypes))
			for _, service := range result.FailingTypes {
				fmt.Printf("  - %s in %s (struct: %v)\n", service.Name, service.Package, service.IsStruct)
			}
		}
	})

	// Test 7: Test namespace filtering
	t.Run("Find all types in models package", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("models").
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d types in models package:\n", len(result.FailingTypes))
			for _, model := range result.FailingTypes {
				fmt.Printf("  - %s\n", model.Name)
			}
		}
	})

	// Test 8: Test negative conditions with ShouldNot
	t.Run("Ensure models don't depend on HTTP", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("models").
			ShouldNot().
			HaveDependencyOn("net/http").
			GetResult()

		if result.IsSuccessful {
			t.Log("✅ Models don't depend on HTTP - good separation!")
		} else {
			t.Errorf("❌ Some models depend on HTTP:")
			for _, model := range result.FailingTypes {
				t.Logf("  - %s has HTTP dependency", model.Name)
			}
		}
	})
}

// TestCreateCustomPredicate demonstrates how to create a custom predicate
func TestCreateCustomPredicate(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Example: Create a custom predicate to find types that contain "Manager" in their name
	t.Run("Custom predicate: Find Manager types", func(t *testing.T) {
		// Get all types first
		allTypes := types.That().GetAllTypes()

		// Apply custom logic
		var managerTypes []*goarchtest.TypeInfo
		for _, typeInfo := range allTypes {
			if strings.Contains(typeInfo.Name, "Manager") {
				managerTypes = append(managerTypes, typeInfo)
			}
		}

		fmt.Printf("Found %d manager types using custom predicate:\n", len(managerTypes))
		for _, manager := range managerTypes {
			fmt.Printf("  - %s in %s\n", manager.Name, manager.Package)
		}

		// You could also test this with the existing NameMatch predicate
		result := types.That().
			HaveNameMatching(".*Manager.*").
			GetResult()

		if result.IsSuccessful {
			fmt.Printf("Found %d manager types using HaveNameMatching:\n", len(result.FailingTypes))
			for _, manager := range result.FailingTypes {
				fmt.Printf("  - %s in %s\n", manager.Name, manager.Package)
			}
		}
	})
}

// TestNamingConventions demonstrates testing naming conventions
func TestNamingConventions(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test that all handlers follow naming convention
	t.Run("Handlers should end with Handler", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("handlers").
			Should().
			HaveNameEndingWith("Handler").
			Or(types.That().ResideInNamespace("handlers").Should().HaveNameEndingWith("Controller")).
			GetResult()

		if result.IsSuccessful {
			t.Log("✅ All handlers follow naming convention")
		} else {
			t.Log("⚠️  Some handlers don't follow naming convention")
		}
	})

	// Test that services are in the correct package
	t.Run("Services should be in services package", func(t *testing.T) {
		result := types.That().
			HaveNameEndingWith("Service").
			Should().
			ResideInNamespace("services").
			GetResult()

		if result.IsSuccessful {
			t.Log("✅ All services are in the correct package")
		} else {
			t.Log("⚠️  Some services are not in the services package")
		}
	})

	// Test that models don't have business logic dependencies
	t.Run("Models should be pure data structures", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("models").
			ShouldNot().
			HaveDependencyOn("services").
			GetResult()

		if result.IsSuccessful {
			t.Log("✅ Models are pure - no service dependencies")
		} else {
			t.Errorf("❌ Some models depend on services:")
			for _, model := range result.FailingTypes {
				t.Logf("  - %s", model.Name)
			}
		}
	})
}
