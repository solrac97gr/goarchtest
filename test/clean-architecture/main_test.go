package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestCleanArchitecture(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Create a clean architecture pattern validator
	// The namespaces are based on our package structure
	cleanArch := goarchtest.CleanArchitecture(
		"github.com/solrac97gr/goarchtest/test/clean_architecture/domain",
		"github.com/solrac97gr/goarchtest/test/clean_architecture/application",
		"github.com/solrac97gr/goarchtest/test/clean_architecture/infrastructure",
		"github.com/solrac97gr/goarchtest/test/clean_architecture/presentation",
	)

	// Validate the architecture
	validationResults := cleanArch.Validate(goarchtest.InPath(projectPath))

	// Debug: Print all packages first
	types := goarchtest.InPath(projectPath)
	allTypes := types.That().GetAllTypes()
	fmt.Println("All types found:")
	for _, t := range allTypes {
		fmt.Printf("- %s in package %s\n", t.Name, t.Package)
		fmt.Printf("  - Full path: %s\n", t.FullPath)
		fmt.Printf("  - Imports: %v\n", t.Imports)
	}

	// Check if any validation rules failed
	for _, result := range validationResults {
		if !result.IsSuccessful {
			t.Errorf("Clean Architecture rule failed: %s", result.RuleDescription)
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
				t.Logf("  - Full path: %s", failingType.FullPath)
				t.Logf("  - Imports: %v", failingType.Imports)
				t.Logf("  - Is Struct: %v", failingType.IsStruct)
				if len(failingType.Interfaces) > 0 {
					t.Logf("  - Interfaces: %v", failingType.Interfaces)
				}
			}
		}
	}
}
