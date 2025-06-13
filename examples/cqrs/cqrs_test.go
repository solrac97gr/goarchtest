package cqrs_test

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestCQRSArchitecture demonstrates testing CQRS (Command Query Responsibility Segregation) architecture
// This test validates that commands and queries are properly separated
func TestCQRSArchitecture(t *testing.T) {
	// Get the absolute path of the project
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	// Create a new Types instance for the project
	types := goarchtest.InPath(projectPath)

	t.Run("Basic CQRS Pattern", func(t *testing.T) {
		// Test basic CQRS pattern
		// Define your CQRS namespaces based on your project structure
		cqrsPattern := goarchtest.CQRSArchitecture(
			"commands",   // Command namespace
			"queries",    // Query namespace
			"domain",     // Domain namespace (shared)
			"writemodel", // Write model namespace
			"readmodel",  // Read model namespace
		)

		// Validate the CQRS pattern
		validationResults := cqrsPattern.Validate(types)

		// Check results
		hasViolations := false
		for i, result := range validationResults {
			if !result.IsSuccessful {
				hasViolations = true
				t.Logf("❌ CQRS Rule #%d failed:", i+1)
				for _, failingType := range result.FailingTypes {
					t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
				}
			} else {
				t.Logf("✅ CQRS Rule #%d passed", i+1)
			}
		}

		if !hasViolations {
			t.Log("🎉 All CQRS architecture rules passed!")
		}
	})

	t.Run("Event Sourced CQRS Pattern", func(t *testing.T) {
		// Test Event Sourced CQRS pattern
		// This is more complex and includes event sourcing concepts
		eventSourcedPattern := goarchtest.EventSourcedCQRSArchitecture(
			"commands",    // Command namespace
			"queries",     // Query namespace
			"events",      // Event namespace
			"eventstore",  // Event store namespace
			"projections", // Projection namespace
			"domain",      // Domain namespace
		)

		// Validate the Event Sourced CQRS pattern
		validationResults := eventSourcedPattern.Validate(types)

		// Check results
		hasViolations := false
		for i, result := range validationResults {
			if !result.IsSuccessful {
				hasViolations = true
				t.Logf("❌ Event Sourced CQRS Rule #%d failed:", i+1)
				for _, failingType := range result.FailingTypes {
					t.Logf("  - %s in package %s", failingType.Name, failingType.Package)
				}
			} else {
				t.Logf("✅ Event Sourced CQRS Rule #%d passed", i+1)
			}
		}

		if !hasViolations {
			t.Log("🎉 All Event Sourced CQRS architecture rules passed!")
		}
	})
}

// TestIndividualCQRSRules demonstrates testing individual CQRS rules
func TestIndividualCQRSRules(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	t.Run("Commands Should Not Depend On Queries", func(t *testing.T) {
		// Commands should not depend on queries (separation of concerns)
		result := types.
			That().
			ResideInNamespace("commands").
			ShouldNot().
			HaveDependencyOn("queries").
			GetResult()

		if !result.IsSuccessful {
			t.Error("❌ CQRS violation: Commands should not depend on queries")
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s depends on queries", failingType.Name)
			}
		} else {
			t.Log("✅ Commands are properly separated from queries")
		}
	})

	t.Run("Queries Should Not Depend On Commands", func(t *testing.T) {
		// Queries should not depend on commands (separation of concerns)
		result := types.
			That().
			ResideInNamespace("queries").
			ShouldNot().
			HaveDependencyOn("commands").
			GetResult()

		if !result.IsSuccessful {
			t.Error("❌ CQRS violation: Queries should not depend on commands")
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s depends on commands", failingType.Name)
			}
		} else {
			t.Log("✅ Queries are properly separated from commands")
		}
	})

	t.Run("Write Models Should Not Depend On Read Models", func(t *testing.T) {
		// Write models should not depend on read models
		result := types.
			That().
			ResideInNamespace("writemodel").
			ShouldNot().
			HaveDependencyOn("readmodel").
			GetResult()

		if !result.IsSuccessful {
			t.Error("❌ CQRS violation: Write models should not depend on read models")
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s depends on read models", failingType.Name)
			}
		} else {
			t.Log("✅ Write models are properly separated from read models")
		}
	})

	t.Run("Commands Should Use Write Models", func(t *testing.T) {
		// Commands should primarily use write models, not read models
		result := types.
			That().
			ResideInNamespace("commands").
			ShouldNot().
			HaveDependencyOn("readmodel").
			GetResult()

		if !result.IsSuccessful {
			t.Error("❌ CQRS violation: Commands should not depend on read models")
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s uses read models instead of write models", failingType.Name)
			}
		} else {
			t.Log("✅ Commands properly use write models")
		}
	})

	t.Run("Queries Should Use Read Models", func(t *testing.T) {
		// Queries should primarily use read models, not write models
		result := types.
			That().
			ResideInNamespace("queries").
			ShouldNot().
			HaveDependencyOn("writemodel").
			GetResult()

		if !result.IsSuccessful {
			t.Error("❌ CQRS violation: Queries should not depend on write models")
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s uses write models instead of read models", failingType.Name)
			}
		} else {
			t.Log("✅ Queries properly use read models")
		}
	})
}

// TestEventSourcingRules demonstrates testing Event Sourcing specific rules
func TestEventSourcingRules(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	t.Run("Commands Should Depend On Events", func(t *testing.T) {
		// Commands should depend on events to produce them
		result := types.
			That().
			ResideInNamespace("commands").
			Should().
			HaveDependencyOn("events").
			GetResult()

		if !result.IsSuccessful {
			t.Log("ℹ️ No event dependencies found in commands (this may be expected if no events are defined)")
		} else {
			t.Log("✅ Commands properly depend on events")
		}
	})

	t.Run("Commands Should Use Event Store", func(t *testing.T) {
		// Commands should interact with event store
		result := types.
			That().
			ResideInNamespace("commands").
			Should().
			HaveDependencyOn("eventstore").
			GetResult()

		if !result.IsSuccessful {
			t.Log("ℹ️ No event store dependencies found in commands (this may be expected if no event store is defined)")
		} else {
			t.Log("✅ Commands properly use event store")
		}
	})

	t.Run("Queries Should Not Use Event Store Directly", func(t *testing.T) {
		// Queries should not depend on event store directly (use projections instead)
		result := types.
			That().
			ResideInNamespace("queries").
			ShouldNot().
			HaveDependencyOn("eventstore").
			GetResult()

		if !result.IsSuccessful {
			t.Error("❌ Event Sourcing violation: Queries should not depend on event store directly")
			for _, failingType := range result.FailingTypes {
				t.Logf("  - %s directly accesses event store", failingType.Name)
			}
		} else {
			t.Log("✅ Queries properly avoid direct event store access")
		}
	})

	t.Run("Projections Should Depend On Events", func(t *testing.T) {
		// Projections should depend on events to build read models
		result := types.
			That().
			ResideInNamespace("projections").
			Should().
			HaveDependencyOn("events").
			GetResult()

		if !result.IsSuccessful {
			t.Log("ℹ️ No event dependencies found in projections (this may be expected if no projections are defined)")
		} else {
			t.Log("✅ Projections properly depend on events")
		}
	})

	t.Run("Queries Should Use Projections", func(t *testing.T) {
		// Queries should depend on projections for read models
		result := types.
			That().
			ResideInNamespace("queries").
			Should().
			HaveDependencyOn("projections").
			GetResult()

		if !result.IsSuccessful {
			t.Log("ℹ️ No projection dependencies found in queries (this may be expected if no projections are defined)")
		} else {
			t.Log("✅ Queries properly use projections")
		}
	})
}
