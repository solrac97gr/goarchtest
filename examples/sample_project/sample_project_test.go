package sample_project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestSampleProject(t *testing.T) {
	// Get the absolute path of the sample project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)

	// Create an error reporter that writes to stderr
	reporter := goarchtest.NewErrorReporter(os.Stderr)

	// Test 1: Domain should not depend on other layers
	result := types.
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("application").
		GetResult()

	reporter.ReportError(result, "Domain should not depend on application layer")

	result = types.
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	reporter.ReportError(result, "Domain should not depend on infrastructure layer")

	result = types.
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("presentation").
		GetResult()

	reporter.ReportError(result, "Domain should not depend on presentation layer")

	// Test 2: Presentation should not directly depend on infrastructure
	// This enforces that presentation goes through application layer
	result = types.
		That().
		ResideInNamespace("presentation").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	reporter.ReportError(result, "Presentation should not directly depend on infrastructure (should go through application)")

	// Test 3: Application can depend on domain (this is correct)
	result = types.
		That().
		ResideInNamespace("application").
		Should().
		HaveDependencyOn("domain").
		GetResult()

	if !result.IsSuccessful {
		reporter.ReportError(result, "Application should depend on domain layer")
	}

	// Test 4: Infrastructure can depend on domain (this is correct)
	result = types.
		That().
		ResideInNamespace("infrastructure").
		Should().
		HaveDependencyOn("domain").
		GetResult()

	if !result.IsSuccessful {
		reporter.ReportError(result, "Infrastructure should depend on domain layer for interfaces")
	}

	// Test 2: Application should only depend on domain
	result = types.
		That().
		ResideInNamespace("application").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	reporter.ReportError(result, "Application should not depend on infrastructure layer")

	result = types.
		That().
		ResideInNamespace("application").
		ShouldNot().
		HaveDependencyOn("presentation").
		GetResult()

	reporter.ReportError(result, "Application should not depend on presentation layer")

	// Test 3: Presentation should only depend on application
	result = types.
		That().
		ResideInNamespace("presentation").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	reporter.ReportError(result, "Presentation should not depend on infrastructure layer")

	// Test 4: Validate using clean architecture pattern
	cleanArchPattern := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
	validationResults := cleanArchPattern.Validate(types)
	reporter.ReportPatternValidation(validationResults)

	// Generate a dependency graph
	// Use the GetAllTypes method to get all types
	allTypes := types.That().GetAllTypes()
	err = reporter.SaveDependencyGraph(allTypes, "dependency_graph.dot")
	if err != nil {
		t.Logf("Failed to save dependency graph: %v", err)
	} else {
		t.Log("Dependency graph saved to dependency_graph.dot")
		t.Log("To generate a PNG, run: dot -Tpng dependency_graph.dot -o dependency_graph.png")
	}
}
