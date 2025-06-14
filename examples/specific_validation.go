package examples

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func SpecificValidation(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Example: Checking if controllers use only services, not repositories
	result := goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("controller").
		Or(goarchtest.InPath(projectPath).That().HaveNameEndingWith("Controller")).
		ShouldNot().
		HaveDependencyOn("repository").
		GetResult()

	if result.IsSuccessful {
		fmt.Println("✅ Controllers don't depend directly on repositories")
	} else {
		fmt.Println("❌ Some controllers directly depend on repositories:")
		for _, failingType := range result.FailingTypes {
			fmt.Printf("  - %s in package %s\n", failingType.Name, failingType.Package)
		}
	}
}
