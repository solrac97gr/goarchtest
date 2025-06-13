package goarchtest_test

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestArchitecture(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Example 1: Types in the presentation package should not have dependencies on data package
	result := goarchtest.InPath(projectPath).
		That().
		ResideInNamespace("presentation").
		ShouldNot().
		HaveDependencyOn("data").
		GetResult()

	if !result.IsSuccessful {
		t.Errorf("Architecture test failed: presentation package has dependencies on data package")
		for _, failingType := range result.FailingTypes {
			t.Logf("Failing type: %s in package %s", failingType.Name, failingType.Package)
		}
	}

	// Example 2: Types that have dependencies on database packages should reside in the data namespace
	result = goarchtest.InPath(projectPath).
		That().
		HaveDependencyOn("database/sql").
		Should().
		ResideInNamespace("data").
		GetResult()

	if !result.IsSuccessful {
		t.Errorf("Architecture test failed: types with database dependencies should be in data package")
		for _, failingType := range result.FailingTypes {
			t.Logf("Failing type: %s in package %s", failingType.Name, failingType.Package)
		}
	}

	// Example 3: All types that implement Repository interface should be structs
	result = goarchtest.InPath(projectPath).
		That().
		ImplementInterface("Repository").
		Should().
		BeStruct().
		GetResult()

	if !result.IsSuccessful {
		t.Errorf("Architecture test failed: Repository implementers should be structs")
		for _, failingType := range result.FailingTypes {
			t.Logf("Failing type: %s in package %s", failingType.Name, failingType.Package)
		}
	}
}
