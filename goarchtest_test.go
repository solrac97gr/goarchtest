package goarchtest_test

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestNew tests the New function which creates a new GoArchTest instance
func TestNew(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	archTest := goarchtest.New(projectPath)

	// Verify the instance is created properly
	if archTest == nil {
		t.Error("New() returned nil")
	}

	// Verify Types is initialized
	if archTest.Types == nil {
		t.Error("New() returned an instance with nil Types")
	}
}

// TestCheckRule tests the CheckRule method
func TestCheckRule(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	archTest := goarchtest.New(projectPath)

	// Test a rule that should pass
	// Types in domain should not depend on infrastructure
	result := archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("domain").
			And().
			DoNotHaveDependencyOn("infrastructure").
			GetResult()
	})

	// Check the result
	if !result.IsSuccessful {
		t.Error("Expected rule to pass, but it failed")
	}

	// Test a rule that should fail
	// We know from the sample project that presentation has violations
	result = archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("presentation").
			And().
			DoNotHaveDependencyOn("infrastructure").
			GetResult()
	})

	// Check the result (should fail due to the intentional violation in the sample project)
	if result.IsSuccessful {
		t.Error("Expected rule to fail, but it passed")
	}
}
