package goarchtest_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// Example demonstrates basic usage of GoArchTest for architectural validation.
func Example() {
	// Get project path
	projectPath, _ := filepath.Abs("./")

	// Test that domain layer doesn't depend on infrastructure
	result := goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("✅ Domain layer is properly isolated")
	} else {
		fmt.Printf("❌ Found %d violations\n", len(result.FailingTypes))
	}
	// Output: ✅ Domain layer is properly isolated
}

// ExampleCleanArchitecture demonstrates using predefined architecture patterns.
func ExampleCleanArchitecture() {
	projectPath, _ := filepath.Abs("./")
	types := goarchtest.InPath(projectPath)

	// Define Clean Architecture pattern
	cleanArch := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")

	// Validate the pattern
	results := cleanArch.Validate(types)

	// Check results
	allPassed := true
	for _, result := range results {
		if !result.IsSuccessful {
			allPassed = false
			break
		}
	}

	if allPassed {
		fmt.Println("✅ Clean Architecture constraints satisfied")
	} else {
		fmt.Println("❌ Clean Architecture violations found")
	}
	// Output: ✅ Clean Architecture constraints satisfied
}

// ExampleTypeSet_HaveNameEndingWith demonstrates filtering types by name suffix.
func ExampleTypeSet_HaveNameEndingWith() {
	projectPath, _ := filepath.Abs("./")

	// Find all services (types ending with "Service")
	result := goarchtest.InPath(projectPath).
		That().
		HaveNameEndingWith("Service").
		Should().
		ResideInNamespace("services").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("✅ All services are in the services package")
	} else {
		fmt.Println("❌ Some services are misplaced")
	}
	// Output: ✅ All services are in the services package
}

// ExampleTypeSet_BeStruct demonstrates filtering for struct types.
func ExampleTypeSet_BeStruct() {
	projectPath, _ := filepath.Abs("./")

	// Find all struct types
	result := goarchtest.InPath(projectPath).
		That().
		BeStruct().
		GetResult()

	fmt.Printf("Found %d struct types\n", len(result.FailingTypes))
	// Output: Found 0 struct types
}

// ExampleTypeSet_AreInterfaces demonstrates filtering for interface types.
func ExampleTypeSet_AreInterfaces() {
	projectPath, _ := filepath.Abs("./")

	// Ensure interfaces are in ports package
	result := goarchtest.InPath(projectPath).
		That().
		AreInterfaces().
		Should().
		ResideInNamespace("ports").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("✅ All interfaces are properly located")
	} else {
		fmt.Println("ℹ️ No interfaces found or some are misplaced")
	}
	// Output: ℹ️ No interfaces found or some are misplaced
}

// ExampleTypeSet_HaveDependencyOn demonstrates dependency analysis.
func ExampleTypeSet_HaveDependencyOn() {
	projectPath, _ := filepath.Abs("./")

	// Ensure types with database dependencies are in data layer
	result := goarchtest.InPath(projectPath).
		That().
		HaveDependencyOn("database/sql").
		Should().
		ResideInNamespace("data").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("✅ Database dependencies are properly isolated")
	} else {
		fmt.Println("ℹ️ No database dependencies found or they're in wrong layer")
	}
	// Output: ℹ️ No database dependencies found or they're in wrong layer
}

// TestArchitecturalConstraints demonstrates comprehensive architectural testing.
func TestArchitecturalConstraints(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test 1: Domain isolation
	t.Run("Domain should not depend on infrastructure", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("domain").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Errorf("Domain layer has infrastructure dependencies:")
			for _, failing := range result.FailingTypes {
				t.Logf("  - %s in %s", failing.Name, failing.Package)
			}
		}
	})

	// Test 2: Service naming conventions
	t.Run("Services should end with Service", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("services").
			Should().
			HaveNameEndingWith("Service").
			GetResult()

		if !result.IsSuccessful {
			t.Log("Some types in services package don't follow naming convention")
		}
	})

	// Test 3: Interface placement
	t.Run("Interfaces should be in ports package", func(t *testing.T) {
		result := types.That().
			AreInterfaces().
			Should().
			ResideInNamespace("ports").
			GetResult()

		if !result.IsSuccessful {
			t.Log("Some interfaces are not in ports package")
		}
	})
}
