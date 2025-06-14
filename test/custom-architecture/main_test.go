package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestHexagonalArchitecture(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	reporter := goarchtest.NewErrorReporter(nil) // nil uses stdout

	// Define the Hexagonal Architecture pattern
	hexagonalPattern := &goarchtest.ArchitecturePattern{
		Name: "Hexagonal Architecture (Ports and Adapters)",
		Rules: []goarchtest.Rule{
			// Core domain should not depend on any external layers
			{
				Description: "Domain layer should not depend on ports",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("core/domain").
						ShouldNot().
						HaveDependencyOn("core/ports").
						GetResult()
				},
			},
			{
				Description: "Domain layer should not depend on adapters",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("core/domain").
						ShouldNot().
						HaveDependencyOn("adapters").
						GetResult()
				},
			},
			// Ports layer should only depend on domain
			{
				Description: "Ports layer should depend on domain",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("core/ports").
						Should().
						HaveDependencyOn("core/domain").
						GetResult()
				},
			},
			{
				Description: "Ports layer should not depend on adapters",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("core/ports").
						ShouldNot().
						HaveDependencyOn("adapters").
						GetResult()
				},
			},
			// Primary adapters should depend on ports but not on secondary adapters
			{
				Description: "Primary adapters should depend on ports",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("adapters/primary").
						Should().
						HaveDependencyOn("core/ports").
						GetResult()
				},
			},
			{
				Description: "Primary adapters should not depend on secondary adapters",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("adapters/primary").
						ShouldNot().
						HaveDependencyOn("adapters/secondary").
						GetResult()
				},
			},
			// Secondary adapters should implement ports but not depend on primary adapters
			{
				Description: "Secondary adapters should depend on ports",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("adapters/secondary").
						Should().
						HaveDependencyOn("core/ports").
						GetResult()
				},
			},
			{
				Description: "Secondary adapters should not depend on primary adapters",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("adapters/secondary").
						ShouldNot().
						HaveDependencyOn("adapters/primary").
						GetResult()
				},
			},
			// Secondary adapters can depend on domain for shared types
			{
				Description: "Secondary adapters should depend on domain",
				Validate: func(types *goarchtest.Types) *goarchtest.Result {
					return types.That().
						ResideInNamespace("adapters/secondary").
						Should().
						HaveDependencyOn("core/domain").
						GetResult()
				},
			},
		},
	}

	// Validate the architecture
	validationResults := hexagonalPattern.Validate(goarchtest.InPath(projectPath))

	// Debug: Print all packages first
	types := goarchtest.InPath(projectPath)
	allTypes := types.That().GetAllTypes()
	fmt.Println("All types found:")
	for _, t := range allTypes {
		fmt.Printf("- %s in package %s\n", t.Name, t.Package)
		fmt.Printf("  - Full path: %s\n", t.FullPath)
		fmt.Printf("  - Imports: %v\n", t.Imports)
		fmt.Printf("  - Is Struct: %v, Is Interface: %v\n", t.IsStruct, t.IsInterface)
	}

	// Report validation results
	reporter.ReportPatternValidation(validationResults)

	// Check if any validation rules failed
	hasFailures := false
	for _, result := range validationResults {
		if !result.IsSuccessful {
			hasFailures = true
			t.Errorf("Hexagonal Architecture rule failed: %s", result.RuleDescription)
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
				t.Logf("  - Full path: %s", failingType.FullPath)
				t.Logf("  - Imports: %v", failingType.Imports)
				t.Logf("  - Is Struct: %v, Is Interface: %v", failingType.IsStruct, failingType.IsInterface)
				if len(failingType.Interfaces) > 0 {
					t.Logf("  - Interface Methods: %v", failingType.Interfaces)
				}
			}
		}
	}

	if !hasFailures {
		t.Log("✅ All Hexagonal Architecture rules passed!")
	}
}

// TestSpecificValidations demonstrates more specific validations
func TestSpecificValidations(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test that domain entities don't have external dependencies
	t.Run("Domain entities should be pure", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("core/domain").
			ShouldNot().
			HaveDependencyOn("net/http").
			GetResult()

		if !result.IsSuccessful {
			t.Errorf("Domain entities should not depend on HTTP: %v", result.FailingTypes)
		}
	})

	// Test that HTTP handlers don't directly depend on database
	t.Run("HTTP handlers should not directly depend on database", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("adapters/primary/http").
			ShouldNot().
			HaveDependencyOn("adapters/secondary/database").
			GetResult()

		if !result.IsSuccessful {
			t.Errorf("HTTP handlers should not directly depend on database: %v", result.FailingTypes)
		}
	})

	// Test naming conventions using name pattern matching
	t.Run("Service implementations should end with Service", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("core/ports").
			And().
			HaveNameEndingWith("Service").
			Should().
			BeStruct().
			GetResult()

		if !result.IsSuccessful {
			t.Logf("Some service implementations don't follow naming convention: %v", result.FailingTypes)
		}
	})

	// Test that repositories end with Repository
	t.Run("Repository implementations should end with Repository", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("adapters/secondary/database").
			Should().
			HaveNameEndingWith("Repository").
			GetResult()

		if !result.IsSuccessful {
			t.Errorf("Repository implementations should end with 'Repository': %v", result.FailingTypes)
		}
	})

	// Test that interfaces are properly defined
	t.Run("Service interfaces should be in ports package", func(t *testing.T) {
		result := types.That().
			AreInterfaces().
			And().
			HaveNameMatching(".*Service").
			Should().
			ResideInNamespace("core/ports").
			GetResult()

		if !result.IsSuccessful {
			t.Logf("Some service interfaces are not in the correct package: %v", result.FailingTypes)
		}
	})

	// Test that all repository interfaces are in ports
	t.Run("Repository interfaces should be in ports package", func(t *testing.T) {
		result := types.That().
			AreInterfaces().
			And().
			HaveNameEndingWith("Repository").
			Should().
			ResideInNamespace("core/ports").
			GetResult()

		if !result.IsSuccessful {
			t.Logf("Some repository interfaces are not in the correct package: %v", result.FailingTypes)
		}
	})

	// Test that handlers follow naming conventions
	t.Run("HTTP handlers should end with Handler", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("adapters/primary/http").
			Should().
			HaveNameEndingWith("Handler").
			GetResult()

		if !result.IsSuccessful {
			t.Logf("Some HTTP handlers don't follow naming convention: %v", result.FailingTypes)
		}
	})
}
