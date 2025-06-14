package examples

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// ExampleCustomPredicate demonstrates how to create and use custom predicates
func ExampleCustomPredicate(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("your_project_path") // Adjust this path as needed
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	isDataAccessComponent := func(typeInfo *goarchtest.TypeInfo) bool {
		return typeInfo.IsStruct &&
			(strings.HasSuffix(typeInfo.Name, "Repository") ||
				strings.HasSuffix(typeInfo.Name, "DAO") ||
				strings.HasSuffix(typeInfo.Name, "Store"))
	}

	result := goarchtest.InPath(projectPath).
		That().
		WithCustomPredicate("IsDataAccessComponent", isDataAccessComponent).
		Should().
		HaveDependencyOn("database/sql").
		Or(goarchtest.InPath(projectPath).That().WithCustomPredicate("IsDataAccessComponent", isDataAccessComponent).HaveDependencyOn("gorm.io")).
		GetResult()

	if result.IsSuccessful {
		fmt.Println("✅ Data access components use appropriate database libraries")
	} else {
		fmt.Println("❌ Some data access components don't use appropriate database libraries:")
		for _, failingType := range result.FailingTypes {
			fmt.Printf("  - %s in package %s\n", failingType.Name, failingType.Package)
		}
	}
}
