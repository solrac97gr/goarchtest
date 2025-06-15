package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestInfrastructureDependencyDetection(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// First, get User domain types
	userDomainTypes := types.That().ResideInNamespace("internal/user/domain").GetAllTypes()
	
	fmt.Printf("=== User domain types and their imports ===\n")
	for _, typeInfo := range userDomainTypes {
		fmt.Printf("Type: %s\n", typeInfo.Name)
		fmt.Printf("  FullPath: %s\n", typeInfo.FullPath)
		fmt.Printf("  Imports: %v\n", typeInfo.Imports)
		
		// Check each import for infrastructure dependency
		hasInfraDep := false
		for _, imp := range typeInfo.Imports {
			if strings.Contains(imp, "infrastructure") {
				hasInfraDep = true
				fmt.Printf("  🚨 INFRASTRUCTURE DEPENDENCY: %s\n", imp)
			}
		}
		if !hasInfraDep {
			fmt.Printf("  ✅ No infrastructure dependencies\n")
		}
		fmt.Printf("---\n")
	}

	// Test HaveDependencyOn with "infrastructure"
	fmt.Printf("\n=== Testing HaveDependencyOn('infrastructure') ===\n")
	infraDepTypes := types.That().
		ResideInNamespace("internal/user/domain").
		HaveDependencyOn("infrastructure").
		GetAllTypes()
	
	fmt.Printf("Found %d types with infrastructure dependency:\n", len(infraDepTypes))
	for _, typeInfo := range infraDepTypes {
		fmt.Printf("  - %s in %s\n", typeInfo.Name, typeInfo.FullPath)
		fmt.Printf("    Imports: %v\n", typeInfo.Imports)
	}

	// Test with more specific path
	fmt.Printf("\n=== Testing HaveDependencyOn('internal/order/infrastructure') ===\n")
	specificInfraDepTypes := types.That().
		ResideInNamespace("internal/user/domain").
		HaveDependencyOn("internal/order/infrastructure").
		GetAllTypes()
	
	fmt.Printf("Found %d types with specific infrastructure dependency:\n", len(specificInfraDepTypes))
	for _, typeInfo := range specificInfraDepTypes {
		fmt.Printf("  - %s in %s\n", typeInfo.Name, typeInfo.FullPath)
		fmt.Printf("    Imports: %v\n", typeInfo.Imports)
	}

	// Test the full violation detection
	fmt.Printf("\n=== Testing ShouldNot + infrastructure dependency ===\n")
	result := types.That().
		ResideInNamespace("internal/user/domain").
		ShouldNot().
		HaveDependencyOn("infrastructure").
		GetResult()
	
	fmt.Printf("Result IsSuccessful: %v\n", result.IsSuccessful)
	fmt.Printf("FailingTypes count: %d\n", len(result.FailingTypes))
	
	if !result.IsSuccessful {
		fmt.Printf("🚨 VIOLATION DETECTED (correctly):\n")
		for _, typeInfo := range result.FailingTypes {
			fmt.Printf("  - %s in %s\n", typeInfo.Name, typeInfo.FullPath)
		}
	} else {
		fmt.Printf("❌ NO VIOLATION DETECTED (this is still the bug!)\n")
	}
}
