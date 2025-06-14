package goarchtest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestReporter tests the Reporter functionality
func TestReporter(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	archTest := goarchtest.New(projectPath)

	// Get some test results
	result1 := archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("domain").
			GetResult()
	})

	result2 := archTest.CheckRule(func(types *goarchtest.Types) *goarchtest.Result {
		return types.That().
			ResideInNamespace("nonexistent").
			GetResult()
	})

	// Create a reporter
	reporter := goarchtest.NewReporter()

	// Add results
	reporter.AddResult(result1)
	reporter.AddResult(result2)

	// Test generating a text report
	reportPath := filepath.Join(os.TempDir(), "goarchtest_report.txt")
	err = reporter.SaveReport("text", reportPath)
	if err != nil {
		t.Errorf("Failed to save text report: %v", err)
	}

	// Check that the report file was created
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Error("Report file was not created")
	} else {
		// Clean up
		os.Remove(reportPath)
	}

	// Test generating an HTML report
	htmlReportPath := filepath.Join(os.TempDir(), "goarchtest_report.html")
	err = reporter.SaveReport("html", htmlReportPath)
	if err != nil {
		t.Errorf("Failed to save HTML report: %v", err)
	}

	// Check that the report file was created
	if _, err := os.Stat(htmlReportPath); os.IsNotExist(err) {
		t.Error("HTML report file was not created")
	} else {
		// Clean up
		os.Remove(htmlReportPath)
	}
}

// TestErrorReporter tests the ErrorReporter functionality
func TestErrorReporter(t *testing.T) {
	// Get the path to the sample project for testing
	projectPath, err := filepath.Abs("./examples/sample_project")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)

	// Create a clean architecture pattern
	cleanArchPattern := goarchtest.CleanArchitecture(
		"domain",
		"application",
		"infrastructure",
		"presentation",
	)

	// Validate the pattern
	results := cleanArchPattern.Validate(types)

	// Create an error reporter with a custom writer
	tempFile, err := os.CreateTemp("", "error_report_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	errorReporter := goarchtest.NewErrorReporter(tempFile)

	// Report pattern validation
	errorReporter.ReportPatternValidation(results)

	// Check that something was written to the file
	fileInfo, err := tempFile.Stat()
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Error report is empty")
	}

	// Test generating a dependency graph
	allTypes := types.That().GetAllTypes()
	dotFilePath := filepath.Join(os.TempDir(), "dependency_graph.dot")

	err = errorReporter.SaveDependencyGraph(allTypes, dotFilePath)
	if err != nil {
		t.Errorf("Failed to save dependency graph: %v", err)
	}

	// Check that the DOT file was created
	if _, err := os.Stat(dotFilePath); os.IsNotExist(err) {
		t.Error("DOT file was not created")
	} else {
		// Clean up
		os.Remove(dotFilePath)
	}
}
