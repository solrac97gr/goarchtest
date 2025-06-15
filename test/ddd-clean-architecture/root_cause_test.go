package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestRootCauseVerification(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)

	// Test 1: Check if we have types initially
	fmt.Printf("=== Initial state ===\n")
	initialTypes := types.That().GetAllTypes()
	fmt.Printf("Initial types count: %d\n", len(initialTypes))
	
	// Test 2: Check after calling That() again
	fmt.Printf("\n=== After second That() call ===\n")
	secondTypes := types.That().GetAllTypes()
	fmt.Printf("Second call types count: %d\n", len(secondTypes))
	
	// Test 3: Try filtering and see what happens
	fmt.Printf("\n=== After filtering ===\n")
	filteredTypes := types.That().ResideInNamespace("internal/user/domain").GetAllTypes()
	fmt.Printf("After ResideInNamespace types count: %d\n", len(filteredTypes))
	
	// Test 4: Check what happens to the original after filtering
	fmt.Printf("\n=== Original after filtering ===\n")
	afterFilterTypes := types.That().GetAllTypes()
	fmt.Printf("Original after filter types count: %d\n", len(afterFilterTypes))
}
