package goarchtest

import (
	"strings"
)

// ResideInNamespace filters types that reside in the specified namespace/package
func (ts *TypeSet) ResideInNamespace(namespace string) *TypeSet {
	ts.currentPredicate = "ResideInNamespace"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if strings.Contains(t.Package, namespace) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// HaveDependencyOn filters types that have a dependency on the specified package
func (ts *TypeSet) HaveDependencyOn(dependency string) *TypeSet {
	ts.currentPredicate = "HaveDependencyOn"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		for _, imp := range t.Imports {
			if strings.Contains(imp, dependency) {
				filteredTypes = append(filteredTypes, t)
				break
			}
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// ImplementInterface filters types that implement the specified interface
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
func (ts *TypeSet) And() *TypeSet {
	ts.currentPredicate = "And"
	// No filtering needed, this is just a logical connector
	return ts
}

// Or performs a union with another TypeSet (logical OR)
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
func (ts *TypeSet) Should() *TypeSet {
	ts.currentPredicate = "Should"
	// Store the current types for later reference
	originalTypes := ts.types
	ts.originalTypes = originalTypes
	return ts
}

// ShouldNot reverses the condition for the following predicates
func (ts *TypeSet) ShouldNot() *TypeSet {
	ts.currentPredicate = "ShouldNot"
	// Store the current types for later reference
	originalTypes := ts.types
	ts.originalTypes = originalTypes
	return ts
}

// Not negates the following predicate
func (ts *TypeSet) Not() *TypeSet {
	ts.currentPredicate = "Not"
	return ts
}
