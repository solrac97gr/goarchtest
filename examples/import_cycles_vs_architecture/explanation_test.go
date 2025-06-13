// Package import_cycles_vs_architecture demonstrates the difference between import cycles and architectural violations
package import_cycles_vs_architecture

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestImportCycleVsArchitecturalViolation demonstrates the key difference between
// what Go's compiler prevents vs what GoArchTest prevents
func TestImportCycleVsArchitecturalViolation(t *testing.T) {
	projectPath, _ := filepath.Abs("../sample_project")
	types := goarchtest.InPath(projectPath)

	t.Run("Go Compiler Prevents Import Cycles", func(t *testing.T) {
		// This would cause a Go compiler error (import cycle):
		//
		// package domain
		// import "infrastructure"
		//
		// package infrastructure
		// import "domain"
		//
		// Go compiler: ERROR - import cycle not allowed

		t.Log("✅ Go compiler automatically prevents A→B→A import cycles")
		t.Log("❌ But Go compiler ALLOWS architectural violations")
	})

	t.Run("GoArchTest Prevents Architectural Violations", func(t *testing.T) {
		// Test 1: Domain should not depend on outer layers (architectural rule)
		result := types.
			That().
			ResideInNamespace("domain").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Logf("✅ GoArchTest caught architectural violation:")
			t.Logf("   Domain layer importing infrastructure (inner→outer dependency)")
			for _, failingType := range result.FailingTypes {
				t.Logf("   Violation: %s depends on infrastructure", failingType.Name)
			}
		}

		// Test 2: Presentation should not skip application layer
		result = types.
			That().
			ResideInNamespace("presentation").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		if !result.IsSuccessful {
			t.Logf("✅ GoArchTest caught layer-skipping violation:")
			t.Logf("   Presentation directly accessing infrastructure (skipping application)")
			for _, failingType := range result.FailingTypes {
				t.Logf("   Violation: %s bypasses application layer", failingType.Name)
			}
		}
	})

	t.Run("What Go Allows But Violates Clean Architecture", func(t *testing.T) {
		examples := []struct {
			description string
			problem     string
			why         string
		}{
			{
				description: "Domain importing infrastructure",
				problem:     "package domain; import \"infrastructure\"",
				why:         "Inner layers should not depend on outer layers",
			},
			{
				description: "Business logic importing HTTP framework",
				problem:     "package business; import \"github.com/gin-gonic/gin\"",
				why:         "Business logic should be framework-agnostic",
			},
			{
				description: "Presentation skipping application layer",
				problem:     "package presentation; import \"infrastructure/repository\"",
				why:         "Should go through application layer for proper separation",
			},
			{
				description: "Cross-domain contamination",
				problem:     "package user; import \"order/internal\"",
				why:         "Bounded contexts should not leak across domains",
			},
		}

		for _, example := range examples {
			t.Logf("❌ %s", example.description)
			t.Logf("   Code: %s", example.problem)
			t.Logf("   Why it's bad: %s", example.why)
			t.Logf("   Go compiler: ✅ Compiles fine")
			t.Logf("   GoArchTest: ❌ Catches violation")
			t.Log("")
		}
	})
}

// TestValueProposition demonstrates the specific value GoArchTest provides
func TestValueProposition(t *testing.T) {
	projectPath, _ := filepath.Abs("../sample_project")
	types := goarchtest.InPath(projectPath)

	t.Run("Architectural Enforcement Beyond Import Cycles", func(t *testing.T) {
		values := []struct {
			capability string
			example    string
		}{
			{
				capability: "Layer Boundary Enforcement",
				example:    "Prevent domain from importing infrastructure",
			},
			{
				capability: "Dependency Direction Control",
				example:    "Ensure dependencies point inward in Clean Architecture",
			},
			{
				capability: "Framework Isolation",
				example:    "Keep business logic free from web framework dependencies",
			},
			{
				capability: "Bounded Context Integrity",
				example:    "Prevent user package from importing order internals",
			},
			{
				capability: "Team Consistency",
				example:    "Automated checks ensure all developers follow architectural rules",
			},
			{
				capability: "Architectural Documentation",
				example:    "Tests serve as executable documentation of architectural decisions",
			},
		}

		for _, value := range values {
			t.Logf("✅ %s: %s", value.capability, value.example)
		}
	})

	t.Run("Clean Architecture Example", func(t *testing.T) {
		// This demonstrates proper Clean Architecture testing
		cleanArchPattern := goarchtest.CleanArchitecture("domain", "application", "infrastructure", "presentation")
		validationResults := cleanArchPattern.Validate(types)

		t.Logf("Clean Architecture Validation Results:")
		for i, result := range validationResults {
			status := "✅ PASS"
			if !result.IsSuccessful {
				status = "❌ FAIL"
			}
			t.Logf("Rule #%d: %s", i+1, status)
		}

		if len(validationResults) > 0 {
			passedRules := 0
			for _, result := range validationResults {
				if result.IsSuccessful {
					passedRules++
				}
			}
			t.Logf("Summary: %d/%d rules passed", passedRules, len(validationResults))
		}
	})
}
