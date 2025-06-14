package examples

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func DefinedArchitecturePattern(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	cleanArch := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
	validationResults := cleanArch.Validate(goarchtest.InPath(projectPath))

	reporter := goarchtest.NewErrorReporter(nil) // nil uses stdout
	reporter.ReportPatternValidation(validationResults)
}
