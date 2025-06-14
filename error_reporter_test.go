package goarchtest_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestErrorReporterOutput tests the output of ErrorReporter functionality
func TestErrorReporterOutput(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	archTest := goarchtest.New(projectPath)

	// Get a failing result
	failResult := archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("nonexistent").
			GetResult()
	})

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create an error reporter with the buffer as writer
	reporter := goarchtest.NewErrorReporter(&buf)

	// Report an error
	reporter.ReportError(failResult, "Test Error Report")

	// Check that the output contains expected text
	output := buf.String()
	if !strings.Contains(output, "Architecture Test Failed: Test Error Report") {
		t.Error("Expected error report to contain failure message")
	}

	if !strings.Contains(output, "Failing Types:") {
		t.Error("Expected error report to list failing types")
	}

	// Reset the buffer
	buf.Reset()

	// Get a successful result
	successResult := archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("domain").
			GetResult()
	})

	// Report a successful result (should not output anything)
	reporter.ReportError(successResult, "Test Success Report")

	// Check that the output is empty
	output = buf.String()
	if output != "" {
		t.Errorf("Expected no output for successful result, got: %s", output)
	}
}

// TestErrorReporterPatternValidation tests the ReportPatternValidation functionality
func TestErrorReporterPatternValidation(t *testing.T) {
	// Get the path to the sample project for testing
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

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create an error reporter with the buffer as writer
	reporter := goarchtest.NewErrorReporter(&buf)

	// Report the pattern validation results
	reporter.ReportPatternValidation(results)

	// Check that the output contains expected text
	output := buf.String()
	if !strings.Contains(output, "Validating Clean Architecture Pattern") {
		t.Error("Expected pattern validation report to mention Clean Architecture")
	}

	if !strings.Contains(output, "Rule:") {
		t.Error("Expected pattern validation report to list rules")
	}
}

// TestErrorReporterDependencyGraph tests the SaveDependencyGraph functionality
func TestErrorReporterDependencyGraph(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new Types instance
	types := goarchtest.InPath(projectPath)

	// Get all types
	allTypes := types.That().GetAllTypes()

	// Create a temporary file path for the graph
	tempDir := t.TempDir()
	graphPath := filepath.Join(tempDir, "test_dependency_graph.dot")

	// Create an error reporter
	reporter := goarchtest.NewErrorReporter(nil)

	// Save dependency graph
	err = reporter.SaveDependencyGraph(allTypes, graphPath)
	if err != nil {
		t.Errorf("Failed to save dependency graph: %v", err)
	}

	// Check that the graph file was created
	if _, err := os.Stat(graphPath); os.IsNotExist(err) {
		t.Error("Graph file was not created")
	}
}
