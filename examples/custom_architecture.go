package examples

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func CustomArchitecturePattern(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	reporter := goarchtest.NewErrorReporter(nil) // nil uses stdout

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
}
