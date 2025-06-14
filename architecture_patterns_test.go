package goarchtest_test

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestCleanArchitecturePattern tests the CleanArchitecture pattern
func TestCleanArchitecturePattern(t *testing.T) {
	// Get the path to the sample project for testing
	// The sample project follows clean architecture
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)

	// Define the clean architecture layers
	domainLayer := "domain"
	applicationLayer := "application"
	infrastructureLayer := "infrastructure"
	presentationLayer := "presentation"

	// Create a clean architecture pattern
	cleanArchPattern := goarchtest.CleanArchitecture(
		domainLayer,
		applicationLayer,
		infrastructureLayer,
		presentationLayer,
	)

	// Validate the pattern
	results := cleanArchPattern.Validate(types)

	// Check the results
	// Note: We expect some rules to fail because the sample project has intentional violations
	// to demonstrate detection capabilities

	// Count successful and failed rules
	successCount := 0
	for _, result := range results {
		if result.IsSuccessful {
			successCount++
		}
	}

	// In clean architecture, we expect a majority of the rules to pass
	// The sample project has a few intentional violations
	if successCount < len(results)/2 {
		t.Errorf("Expected majority of clean architecture rules to pass, but only %d of %d passed",
			successCount, len(results))
	}
}

// TestIndividualCleanArchitectureRules tests individual rules of Clean Architecture
func TestIndividualCleanArchitectureRules(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	archTest := goarchtest.New(projectPath)

	// Test that application depends on domain (should pass)
	result := archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("application").
			Should().
			HaveDependencyOn("domain").
			GetResult()
	})

	if !result.IsSuccessful {
		t.Error("Expected application to depend on domain")
	}

	// Test that infrastructure depends on domain (should pass)
	result = archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("infrastructure").
			Should().
			HaveDependencyOn("domain").
			GetResult()
	})

	if !result.IsSuccessful {
		t.Error("Expected infrastructure to depend on domain")
	}

	// Test that presentation depends on application (should pass)
	result = archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("presentation").
			Should().
			HaveDependencyOn("application").
			GetResult()
	})

	if !result.IsSuccessful {
		t.Error("Expected presentation to depend on application")
	}

	// Test that domain doesn't depend on infrastructure for the clean User (should pass)
	result = archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("domain").
			And().
			NameMatch("^User$"). // Only pure User, not UserWithViolation
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()
	})

	if !result.IsSuccessful {
		t.Error("Expected clean domain User to not depend on infrastructure")
	}
}

// TestDDDWithCleanArchitecture tests the DDD with Clean Architecture pattern
func TestDDDWithCleanArchitecture(t *testing.T) {
	// Get the path to the DDD example project
	projectPath, err := filepath.Abs("./examples/ddd_clean_architecture")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)

	// Define the bounded contexts (domains)
	domains := []string{"user", "products"}
	sharedKernel := "internal/shared"
	utils := "pkg"

	// Create a DDD with Clean Architecture pattern
	dddPattern := goarchtest.DDDWithCleanArchitecture(
		domains,
		sharedKernel,
		utils,
	)

	// Validate the pattern
	results := dddPattern.Validate(types)

	// Count successful and failed rules
	successCount := 0
	for _, result := range results {
		if result.IsSuccessful {
			successCount++
		}
	}

	// In DDD with Clean Architecture, we expect a majority of the rules to pass
	// The example project has a few intentional violations
	if successCount < len(results)/2 {
		t.Errorf("Expected majority of DDD with Clean Architecture rules to pass, but only %d of %d passed",
			successCount, len(results))
	}
}

// TestPatternCombinations tests combining different architectural patterns
func TestPatternCombinations(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	archTest := goarchtest.New(projectPath)
	types := archTest.Types

	// Define a custom set of rules combining aspects of different patterns

	// Rule 1: Domain types should be structs (from Clean Architecture)
	rule1 := func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("domain").
			Should().
			BeStruct().
			GetResult()
	}

	// Rule 2: Repository implementations should end with "Repository" (naming convention)
	rule2 := func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("infrastructure").
			Should().
			HaveNameEndingWith("Repository").
			GetResult()
	}

	// Rule 3: Handlers should depend on application services (from Clean Architecture)
	rule3 := func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("presentation").
			Should().
			HaveDependencyOn("application").
			GetResult()
	}

	// Execute the rules
	result1 := rule1(types)
	result2 := rule2(types)
	result3 := rule3(types)

	// Check the results
	if !result1.IsSuccessful {
		t.Error("Rule 1 failed: Domain types should be structs")
	}

	if !result2.IsSuccessful {
		t.Error("Rule 2 failed: Repository implementations should end with 'Repository'")
	}

	if !result3.IsSuccessful {
		t.Error("Rule 3 failed: Handlers should depend on application services")
	}
}
