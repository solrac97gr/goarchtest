package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func ExampleSimpleArchitectureTest(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	fmt.Println("Running architecture tests...")

	// Example 1: Test that presentation layer doesn't depend on data layer
	result := goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("presentation").
		ShouldNot().
		HaveDependencyOn("data").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("Test 1 passed: Presentation layer doesn't depend on data layer")
	} else {
		fmt.Println("Test 1 failed: Presentation layer has dependencies on data layer")
		for _, failingType := range result.FailingTypes {
			fmt.Printf("  - %s in package %s\n", failingType.Name, failingType.Package)
		}
	}

	// Example 2: Test that types with database dependencies are in data package
	result = goarchtest.InPath(projectPath).
		That().
		HaveDependencyOn("database/sql").
		Should().
		ResideInNamespace("data").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("Test 2 passed: Types with database dependencies are in data package")
	} else {
		fmt.Println("Test 2 failed: Types with database dependencies should be in data package")
		for _, failingType := range result.FailingTypes {
			fmt.Printf("  - %s in package %s\n", failingType.Name, failingType.Package)
		}
	}

	// Example 3: Test naming conventions
	result = goarchtest.InPath(projectPath).
		That().
		ImplementInterface("Repository").
		Should().
		HaveNameEndingWith("Repository").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("Test 3 passed: Repository implementations follow naming convention")
	} else {
		fmt.Println("Test 3 failed: Repository implementations should end with 'Repository'")
		for _, failingType := range result.FailingTypes {
			fmt.Printf("  - %s in package %s\n", failingType.Name, failingType.Package)
		}
	}
}
