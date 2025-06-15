package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestDebugDependencyDetection tests if we can detect the architectural violation
func TestDebugDependencyDetection(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Debug: Print all types and their dependencies
	allTypes := types.That().GetAllTypes()
	fmt.Printf("=== DEBUG: All types and their imports ===\n")
	for _, typeInfo := range allTypes {
		fmt.Printf("Type: %s\n", typeInfo.Name)
		fmt.Printf("  Package: %s\n", typeInfo.Package)
		fmt.Printf("  FullPath: %s\n", typeInfo.FullPath)
		fmt.Printf("  Imports: %v\n", typeInfo.Imports)
		fmt.Printf("  IsStruct: %v, IsInterface: %v\n", typeInfo.IsStruct, typeInfo.IsInterface)
		fmt.Printf("---\n")
	}
	fmt.Printf("=== END DEBUG ===\n\n")

	// Test specific violation: User domain importing infrastructure
	t.Run("Debug - User domain types", func(t *testing.T) {
		userDomainTypes := types.That().
			ResideInNamespace("internal/user/domain").
			GetResult()

		fmt.Printf("=== User domain types ===\n")
		for _, typeInfo := range userDomainTypes.FailingTypes {
			fmt.Printf("Type: %s in %s\n", typeInfo.Name, typeInfo.FullPath)
			fmt.Printf("  Imports: %v\n", typeInfo.Imports)
			for _, imp := range typeInfo.Imports {
				if contains(imp, "infrastructure") {
					fmt.Printf("  🚨 FOUND INFRASTRUCTURE DEPENDENCY: %s\n", imp)
				}
			}
		}
	})

	// Test with HaveDependencyOn
	t.Run("Debug - Types with infrastructure dependency", func(t *testing.T) {
		result := types.That().
			HaveDependencyOn("infrastructure").
			GetResult()

		fmt.Printf("=== Types with infrastructure dependency ===\n")
		fmt.Printf("Found %d types with infrastructure dependency\n", len(result.FailingTypes))
		for _, typeInfo := range result.FailingTypes {
			fmt.Printf("Type: %s in %s\n", typeInfo.Name, typeInfo.FullPath)
			fmt.Printf("  Imports: %v\n", typeInfo.Imports)
		}
	})

	// Test the actual rule that should fail
	t.Run("Debug - Domain should not depend on infrastructure", func(t *testing.T) {
		result := types.That().
			ResideInNamespace("internal/user/domain").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()

		fmt.Printf("=== Domain dependency test result ===\n")
		fmt.Printf("IsSuccessful: %v\n", result.IsSuccessful)
		fmt.Printf("FailingTypes count: %d\n", len(result.FailingTypes))
		
		if !result.IsSuccessful {
			fmt.Printf("🚨 VIOLATION DETECTED (as expected):\n")
			for _, typeInfo := range result.FailingTypes {
				fmt.Printf("  - %s in %s\n", typeInfo.Name, typeInfo.FullPath)
			}
		} else {
			fmt.Printf("❌ NO VIOLATION DETECTED (this is the bug!)\n")
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || 
		   len(s) > len(substr) && s[:len(substr)] == substr ||
		   len(s) > len(substr) && s[len(s)-len(substr)-1:len(s)-1] == substr
}
