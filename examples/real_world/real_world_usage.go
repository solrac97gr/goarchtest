package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func ExampleRealWorldUsage(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	fmt.Println("Running real-world architecture tests...")

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

	// Example: Using predefined architecture patterns
	cleanArch := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
	validationResults := cleanArch.Validate(goarchtest.InPath(projectPath))

	reporter := goarchtest.NewErrorReporter(nil) // nil uses stdout
	reporter.ReportPatternValidation(validationResults)

	// Example: Creating a custom architecture pattern
	microservicesPattern := &goarchtest.ArchitecturePattern{
		Name: "Microservices Pattern",
		Rules: []func(*goarchtest.Types) *goarchtest.Result{
			// Services should be independent from each other
			func(types *goarchtest.Types) *goarchtest.Result {
				return types.That().
					ResideInNamespace("services/user").
					ShouldNot().
					HaveDependencyOn("services/product").
					GetResult()
			},
			func(types *goarchtest.Types) *goarchtest.Result {
				return types.That().
					ResideInNamespace("services/product").
					ShouldNot().
					HaveDependencyOn("services/user").
					GetResult()
			},
			// Services should communicate via API clients
			func(types *goarchtest.Types) *goarchtest.Result {
				return types.That().
					ResideInNamespace("services").
					And().
					HaveDependencyOn("services").
					Should().
					HaveDependencyOn("api/client").
					GetResult()
			},
		},
	}

	microservicesResults := microservicesPattern.Validate(goarchtest.InPath(projectPath))
	reporter.ReportPatternValidation(microservicesResults)

	// Example: Custom predicate for unique architecture rules
	isDataAccessComponent := func(typeInfo *goarchtest.TypeInfo) bool {
		return typeInfo.IsStruct &&
			(strings.HasSuffix(typeInfo.Name, "Repository") ||
				strings.HasSuffix(typeInfo.Name, "DAO") ||
				strings.HasSuffix(typeInfo.Name, "Store"))
	}

	result = goarchtest.InPath(projectPath).
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
