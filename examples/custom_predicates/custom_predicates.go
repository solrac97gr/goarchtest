package custom_predicates

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// ExampleCustomPredicate demonstrates how to create and use custom predicates
func ExampleCustomPredicate(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)

	// Define a custom predicate function
	isServiceImplementation := func(typeInfo *goarchtest.TypeInfo) bool {
		// A service implementation should:
		// 1. End with "Service"
		// 2. Be a struct
		return typeInfo.IsStruct && len(typeInfo.Name) > 7 && typeInfo.Name[len(typeInfo.Name)-7:] == "Service"
	}

	// Use our custom predicate with the helper function
	result := types.
		That().
		WithCustomPredicate("IsServiceImplementation", isServiceImplementation).
		Should().
		ResideInNamespace("application").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Service implementations should reside in the application layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}
}
