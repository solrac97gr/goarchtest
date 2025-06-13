package clean_architecture

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestCleanArchitecture validates that the project follows clean architecture principles
func TestCleanArchitecture(t *testing.T) {
	// Get the absolute path of the project to test
	projectPath, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// 1. Domain layer should not depend on any other layer
	result := goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Domain layer should not depend on infrastructure layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("application").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Domain layer should not depend on application layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("domain").
		ShouldNot().
		HaveDependencyOn("presentation").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Domain layer should not depend on presentation layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	// 2. Application layer should only depend on domain layer
	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("application").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Application layer should not depend on infrastructure layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("application").
		ShouldNot().
		HaveDependencyOn("presentation").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Application layer should not depend on presentation layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	// 3. Presentation layer should only depend on application layer
	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("presentation").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Presentation layer should not depend on infrastructure layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	// 4. Infrastructure layer should not depend on presentation layer
	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("infrastructure").
		ShouldNot().
		HaveDependencyOn("presentation").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Infrastructure layer should not depend on presentation layer")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	// 5. Check naming conventions
	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("application").
		And().
		ImplementInterface("Service").
		Should().
		HaveNameEndingWith("Service").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Service implementations should end with 'Service'")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}

	result = goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("infrastructure").
		And().
		ImplementInterface("Repository").
		Should().
		HaveNameEndingWith("Repository").
		GetResult()

	if !result.IsSuccessful {
		t.Error("Repository implementations should end with 'Repository'")
		for _, failingType := range result.FailingTypes {
			t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
		}
	}
}
