package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// ExampleErrorHandling demonstrates how to use advanced error handling features
func ExampleErrorHandling(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)

	// Test architecture rule
	result := types.That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	// Advanced error handling with detailed failure information
	if !result.IsSuccessful {
		// Get detailed failure information
		details := result.GetFailureDetails()

		// Log the details
		t.Errorf("Architecture violation: Domain depends on infrastructure\n%s", details)

		// You could also use this for custom error reporting or formatting
		fmt.Printf("Error report:\n%s\n", details)

		// Or in a CI environment, you might want to output it in a specific format
		fmt.Printf("::error::Architecture test failed with %d violations\n", len(result.FailingTypes))
	}

	// Using with ErrorReporter for more sophisticated reporting
	reporter := goarchtest.NewErrorReporter(nil) // nil uses stdout

	// Report error with custom description
	reporter.ReportError(result, "Domain layer must not depend on infrastructure")

	// Test multiple architectural patterns
	cleanArch := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
	validationResults := cleanArch.Validate(types)

	// Report pattern validation results
	reporter.ReportPatternValidation(validationResults)
}
