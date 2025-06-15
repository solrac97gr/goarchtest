package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestSpecificViolationDetection(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test 1: Check ResideInNamespace filtering
	t.Run("ResideInNamespace filtering", func(t *testing.T) {
		userDomainTypes := types.That().
			ResideInNamespace("internal/user/domain").
			GetResult()

		fmt.Printf("=== ResideInNamespace('internal/user/domain') results ===\n")
		fmt.Printf("IsSuccessful: %v\n", userDomainTypes.IsSuccessful)
		fmt.Printf("FailingTypes count: %d\n", len(userDomainTypes.FailingTypes))
		
		allTypes := types.That().GetAllTypes()
		fmt.Printf("Total types in project: %d\n", len(allTypes))
		
		// Count types that should match
		expectedMatches := 0
		for _, typeInfo := range allTypes {
			if strings.Contains(typeInfo.FullPath, "internal/user/domain") {
				expectedMatches++
				fmt.Printf("Expected match: %s in %s\n", typeInfo.Name, typeInfo.FullPath)
			}
		}
		fmt.Printf("Expected matches: %d\n", expectedMatches)
	})

	// Test 2: Check HaveDependencyOn filtering  
	t.Run("HaveDependencyOn infrastructure filtering", func(t *testing.T) {
		infraDeps := types.That().
			HaveDependencyOn("infrastructure").
			GetResult()

		fmt.Printf("\n=== HaveDependencyOn('infrastructure') results ===\n")
		fmt.Printf("IsSuccessful: %v\n", infraDeps.IsSuccessful)
		fmt.Printf("FailingTypes count: %d\n", len(infraDeps.FailingTypes))
		
		allTypes := types.That().GetAllTypes()
		
		// Find types that actually have infrastructure dependencies
		actualMatches := 0
		for _, typeInfo := range allTypes {
			for _, imp := range typeInfo.Imports {
				if strings.Contains(imp, "infrastructure") {
					actualMatches++
					fmt.Printf("Type %s has infrastructure dependency: %s\n", typeInfo.Name, imp)
					break
				}
			}
		}
		fmt.Printf("Actual types with infrastructure deps: %d\n", actualMatches)
	})

	// Test 3: Manual check of the User type specifically
	t.Run("Manual User type check", func(t *testing.T) {
		allTypes := types.That().GetAllTypes()
		
		for _, typeInfo := range allTypes {
			if typeInfo.Name == "User" && strings.Contains(typeInfo.FullPath, "internal/user/domain/models") {
				fmt.Printf("\n=== User type analysis ===\n")
				fmt.Printf("Name: %s\n", typeInfo.Name)
				fmt.Printf("Package: %s\n", typeInfo.Package)
				fmt.Printf("FullPath: %s\n", typeInfo.FullPath)
				fmt.Printf("Imports: %v\n", typeInfo.Imports)
				
				// Check if it has infrastructure dependency
				hasInfraDep := false
				for _, imp := range typeInfo.Imports {
					if strings.Contains(imp, "infrastructure") {
						hasInfraDep = true
						fmt.Printf("🚨 VIOLATION: User imports infrastructure: %s\n", imp)
					}
				}
				
				if !hasInfraDep {
					fmt.Printf("✅ No infrastructure dependency found\n")
				}
				
				// Check if it should be caught by ResideInNamespace
				shouldMatch := strings.Contains(typeInfo.FullPath, "internal/user/domain")
				fmt.Printf("Should match 'internal/user/domain': %v\n", shouldMatch)
				
				break
			}
		}
	})

	// Test 4: Test the actual problematic query step by step
	t.Run("Step by step violation test", func(t *testing.T) {
		fmt.Printf("\n=== Step by step test ===\n")
		
		// Step 1: Get all types
		allTypes := types.That().GetAllTypes()
		fmt.Printf("1. Total types: %d\n", len(allTypes))
		
		// Step 2: Filter by namespace
		step2 := types.That().ResideInNamespace("internal/user/domain")
		step2Types := step2.GetAllTypes()
		fmt.Printf("2. After ResideInNamespace('internal/user/domain'): %d types\n", len(step2Types))
		for _, t := range step2Types {
			fmt.Printf("   - %s in %s\n", t.Name, t.FullPath)
		}
		
		// Step 3: Filter by dependency
		step3 := types.That().
			ResideInNamespace("internal/user/domain").
			HaveDependencyOn("infrastructure")
		step3Types := step3.GetAllTypes()
		fmt.Printf("3. After adding HaveDependencyOn('infrastructure'): %d types\n", len(step3Types))
		for _, t := range step3Types {
			fmt.Printf("   - %s in %s\n", t.Name, t.FullPath)
			fmt.Printf("     Imports: %v\n", t.Imports)
		}
		
		// Step 4: Apply ShouldNot and get result
		result := types.That().
			ResideInNamespace("internal/user/domain").
			ShouldNot().
			HaveDependencyOn("infrastructure").
			GetResult()
		
		fmt.Printf("4. Final result with ShouldNot:\n")
		fmt.Printf("   IsSuccessful: %v\n", result.IsSuccessful)
		fmt.Printf("   FailingTypes count: %d\n", len(result.FailingTypes))
		
		if result.IsSuccessful {
			fmt.Printf("   🚨 BUG: Test passed when it should have failed!\n")
		} else {
			fmt.Printf("   ✅ Test correctly failed\n")
			for _, t := range result.FailingTypes {
				fmt.Printf("   Violating type: %s in %s\n", t.Name, t.FullPath)
			}
		}
	})
}
