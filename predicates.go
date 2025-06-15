package goarchtest

import (
	"strings"
)

// ResideInNamespace filters types that reside in the specified namespace/package
// It allows for filtering based on the package namespace of the type.
// Parameters:
//   - namespace: A string representing the namespace to match against type packages
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose packages match the namespace,
//     allowing for method chaining
//
// Example:
//
//	typeSet.ResideInNamespace("github.com/myorg/mypackage")
func (ts *TypeSet) ResideInNamespace(namespace string) *TypeSet {
	ts.currentPredicate = "ResideInNamespace"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		// Check exact match first
		if t.FullPath == namespace {
			filteredTypes = append(filteredTypes, t)
			continue
		}
		
		// Check if namespace matches the end of the FullPath (relative path matching)
		if strings.HasSuffix(t.FullPath, "/"+namespace) || strings.Contains(t.FullPath, "/"+namespace+"/") {
			filteredTypes = append(filteredTypes, t)
			continue
		}
		
		// Also check prefix match for full paths
		if strings.HasPrefix(t.FullPath, namespace+"/") {
			filteredTypes = append(filteredTypes, t)
			continue
		}
	}

	// Create a new TypeSet to avoid modifying the original
	newTypeSet := &TypeSet{
		types:             filteredTypes,
		originalTypes:     ts.originalTypes, // Keep reference to original types
		currentPredicate:  ts.currentPredicate,
		matchedPredicates: append([]string{}, ts.matchedPredicates...), // Copy slice
	}
	newTypeSet.matchedPredicates = append(newTypeSet.matchedPredicates, ts.currentPredicate)
	return newTypeSet
}

// HaveDependencyOn filters types that have a dependency on the specified package
// It allows for filtering based on the import statements of the type.
// Parameters:
//   - dependency: A string representing the package dependency to match against type imports
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types that import the specified dependency,
//     allowing for method chaining
//
// Example:
//
//	typeSet.HaveDependencyOn("github.com/some/dependency")
func (ts *TypeSet) HaveDependencyOn(dependency string) *TypeSet {
	ts.currentPredicate = "HaveDependencyOn"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		for _, imp := range t.Imports {
			// Exact match
			if imp == dependency {
				filteredTypes = append(filteredTypes, t)
				break
			}
			
			// Prefix match with slash (for exact package boundaries)
			if strings.HasPrefix(imp, dependency+"/") {
				filteredTypes = append(filteredTypes, t)
				break
			}
			
			// Suffix match for relative path matching (e.g., "infrastructure" matches "*/infrastructure")
			if strings.HasSuffix(imp, "/"+dependency) {
				filteredTypes = append(filteredTypes, t)
				break
			}
			
			// Contains match for partial path matching (e.g., "infrastructure" matches "*/infrastructure/*")
			if strings.Contains(imp, "/"+dependency+"/") {
				filteredTypes = append(filteredTypes, t)
				break
			}
		}
	}

	// Create a new TypeSet to avoid modifying the original
	newTypeSet := &TypeSet{
		types:             filteredTypes,
		originalTypes:     ts.originalTypes, // Keep reference to original types
		currentPredicate:  ts.currentPredicate,
		matchedPredicates: append([]string{}, ts.matchedPredicates...), // Copy slice
	}
	newTypeSet.matchedPredicates = append(newTypeSet.matchedPredicates, ts.currentPredicate)
	return newTypeSet
}

// ImplementInterface filters types that implement the specified interface
// It allows for filtering based on the interfaces implemented by the type.
// Parameters:
//   - interfaceName: A string representing the name of the interface to check against
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types that implement the specified interface,
//     allowing for method chaining
//
// Example:
//
//	typeSet.ImplementInterface("MyInterface")
func (ts *TypeSet) ImplementInterface(interfaceName string) *TypeSet {
	ts.currentPredicate = "ImplementInterface"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		for _, iface := range t.Interfaces {
			if iface == interfaceName {
				filteredTypes = append(filteredTypes, t)
				break
			}
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// BeStruct filters types that are structs
// It allows for filtering based on whether the type is a struct.
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only struct types,
//     allowing for method chaining
//
// Example:
//
//	typeSet.BeStruct()
func (ts *TypeSet) BeStruct() *TypeSet {
	ts.currentPredicate = "BeStruct"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if t.IsStruct {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// And combines predicates (logical AND)
// It allows for chaining multiple predicates together, ensuring that all conditions must be met.
// Returns:
//   - *TypeSet: Returns the TypeSet itself to allow for method chaining
//
// Example:
//
//	typeSet.And().HaveDependencyOn("github.com/some/dependency").BeStruct()
func (ts *TypeSet) And() *TypeSet {
	ts.currentPredicate = "And"
	// No filtering needed, this is just a logical connector
	return ts
}

// Or performs a union with another TypeSet (logical OR)
// It allows for combining two TypeSets, resulting in a new TypeSet that contains types from both sets.
// Returns:
//   - *TypeSet: Returns a new TypeSet that is the union of the two sets, allowing for method chaining
//
// Example:
//
//	typeSet1.Or(typeSet2)
func (ts *TypeSet) Or(other *TypeSet) *TypeSet {
	ts.currentPredicate = "Or"

	// Create a union of the two type sets
	unionMap := make(map[string]bool)
	for _, t := range ts.types {
		key := t.Package + "." + t.Name
		unionMap[key] = true
	}

	for _, t := range other.types {
		key := t.Package + "." + t.Name
		if !unionMap[key] {
			ts.types = append(ts.types, t)
			unionMap[key] = true
		}
	}

	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// Should reverses the condition for the following predicates
// It allows for asserting that the following predicates should hold true.
// Returns:
//   - *TypeSet: Returns the TypeSet itself to allow for method chaining
//
// Example:
//
//	ts.Should().HaveDependencyOn("github.com/some/dependency").BeStruct()
func (ts *TypeSet) Should() *TypeSet {
	ts.currentPredicate = "Should"
	// Store the current types for later reference
	originalTypes := ts.types
	ts.originalTypes = originalTypes
	return ts
}

// ShouldNot reverses the condition for the following predicates
// It allows for asserting that the following predicates should not hold true.
// Returns:
//   - *TypeSet: Returns the TypeSet itself to allow for method chaining
//
// Example:
//
//	ts.ShouldNot().HaveDependencyOn("github.com/some/dependency").BeStruct()
func (ts *TypeSet) ShouldNot() *TypeSet {
	ts.currentPredicate = "ShouldNot"
	// Create a new TypeSet to avoid modifying the original
	newTypeSet := &TypeSet{
		types:             append([]*TypeInfo{}, ts.types...), // Copy types slice
		originalTypes:     ts.originalTypes,
		currentPredicate:  ts.currentPredicate,
		matchedPredicates: append([]string{}, ts.matchedPredicates...), // Copy slice
	}
	newTypeSet.matchedPredicates = append(newTypeSet.matchedPredicates, "Negate")
	return newTypeSet
}

// Not negates the following predicate
// It allows for negating the condition of the next predicate.
// Returns:
//   - *TypeSet: Returns the TypeSet itself to allow for method chaining
//
// Example:
//
//	ts.Not().HaveDependencyOn("github.com/some/dependency")
func (ts *TypeSet) Not() *TypeSet {
	ts.currentPredicate = "Not"
	return ts
}
