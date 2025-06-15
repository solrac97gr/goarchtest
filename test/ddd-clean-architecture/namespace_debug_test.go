package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestNamespaceMatching(t *testing.T) {
	projectPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Failed to get project path: %v", err)
	}

	types := goarchtest.InPath(projectPath)
	allTypes := types.That().GetAllTypes()

	fmt.Printf("=== Debugging namespace matching ===\n")
	fmt.Printf("Total types: %d\n\n", len(allTypes))

	searchNamespace := "internal/user/domain"
	fmt.Printf("Searching for namespace: '%s'\n\n", searchNamespace)

	for _, typeInfo := range allTypes {
		fmt.Printf("Type: %s\n", typeInfo.Name)
		fmt.Printf("  FullPath: '%s'\n", typeInfo.FullPath)
		
		// Test exact match
		exactMatch := typeInfo.FullPath == searchNamespace
		fmt.Printf("  Exact match ('%s' == '%s'): %v\n", typeInfo.FullPath, searchNamespace, exactMatch)
		
		// Test prefix match
		prefixMatch := strings.HasPrefix(typeInfo.FullPath, searchNamespace+"/")
		fmt.Printf("  Prefix match (starts with '%s/'): %v\n", searchNamespace, prefixMatch)
		
		// Test contains
		containsMatch := strings.Contains(typeInfo.FullPath, searchNamespace)
		fmt.Printf("  Contains match: %v\n", containsMatch)
		
		shouldMatch := exactMatch || prefixMatch
		fmt.Printf("  Should match: %v\n", shouldMatch)
		fmt.Printf("---\n")
	}

	// Test the actual filtering
	fmt.Printf("\n=== Testing actual ResideInNamespace filtering ===\n")
	filtered := types.That().ResideInNamespace(searchNamespace)
	filteredTypes := filtered.GetAllTypes()
	fmt.Printf("ResideInNamespace('%s') returned %d types\n", searchNamespace, len(filteredTypes))
	
	for _, typeInfo := range filteredTypes {
		fmt.Printf("  - %s in %s\n", typeInfo.Name, typeInfo.FullPath)
	}
}
