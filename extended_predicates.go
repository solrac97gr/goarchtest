package goarchtest

import (
	"regexp"
	"strings"
)

// NameMatch filters types based on a regex pattern match on their names.
// It allows for flexible matching of type names using regular expressions.
//
// Parameters:
//   - pattern: A string representing the regex pattern to match against type names
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose names match the pattern,
//     allowing for method chaining
//
// Example:
//
//	typeSet.NameMatch("MyType.*")
func (ts *TypeSet) NameMatch(pattern string) *TypeSet {
	ts.currentPredicate = "NameMatch"

	regex, err := regexp.Compile(pattern)
	if err != nil {
		// If pattern is invalid, return empty set
		ts.types = []*TypeInfo{}
		return ts
	}

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if regex.MatchString(t.Name) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// HaveNameEndingWith filters types whose names end with the specified suffix.
// It allows for easy identification of types based on their naming conventions.
// Parameters:
//   - suffix: A string representing the suffix to match against type names
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose names end with the suffix,
//     allowing for method chaining
//
// Example:
//
//	typeSet.HaveNameEndingWith("Handler")
func (ts *TypeSet) HaveNameEndingWith(suffix string) *TypeSet {
	ts.currentPredicate = "HaveNameEndingWith"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if strings.HasSuffix(t.Name, suffix) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// HaveNameStartingWith filters types whose names start with the specified prefix
// It allows for easy identification of types based on their naming conventions.
// Parameters:
//   - prefix: A string representing the prefix to match against type names
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose names start with the prefix,
//     allowing for method chaining
//
// Example:
//
//	typeSet.HaveNameStartingWith("My")
func (ts *TypeSet) HaveNameStartingWith(prefix string) *TypeSet {
	ts.currentPredicate = "HaveNameStartingWith"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if strings.HasPrefix(t.Name, prefix) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// ResideInDirectory filters types that reside in the specified directory
// It allows for filtering based on the directory structure of the type's full path.
// Parameters:
//   - directory: A string representing the directory to match against type full paths
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose full paths contain the directory,
//     allowing for method chaining
//
// Example:
//
//	typeSet.ResideInDirectory("internal/mydir")
func (ts *TypeSet) ResideInDirectory(directory string) *TypeSet {
	ts.currentPredicate = "ResideInDirectory"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if strings.Contains(t.FullPath, directory) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// DoNotResideInNamespace filters types that do not reside in the specified namespace
// It allows for excluding types based on their package namespace.
// Parameters:
//   - namespace: A string representing the namespace to exclude from the type's package
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose packages do not contain the namespace,
//     allowing for method chaining
//
// Example:
//
//	typeSet.DoNotResideInNamespace("github.com/external/pkg")
func (ts *TypeSet) DoNotResideInNamespace(namespace string) *TypeSet {
	ts.currentPredicate = "DoNotResideInNamespace"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if !strings.Contains(t.Package, namespace) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// DoNotHaveDependencyOn filters the TypeSet to include only types that do not have
// a dependency on the specified import path. A type is considered to not have the
// dependency if none of its imports contain the given dependency string.
//
// Parameters:
//   - dependency: A string representing the import path or part of it to check against
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types without the specified dependency,
//     allowing for method chaining
//
// Example:
//
//	typeSet.DoNotHaveDependencyOn("github.com/external/pkg")
func (ts *TypeSet) DoNotHaveDependencyOn(dependency string) *TypeSet {
	ts.currentPredicate = "DoNotHaveDependencyOn"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		hasDependency := false
		for _, imp := range t.Imports {
			if strings.Contains(imp, dependency) {
				hasDependency = true
				break
			}
		}

		if !hasDependency {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}

// HaveNameMatching filters types based on a regex pattern match on their names.
// This is an alias for NameMatch for better readability in test scenarios.
//
// Parameters:
//   - pattern: A string representing the regex pattern to match against type names
//
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only types whose names match the pattern,
//     allowing for method chaining
//
// Example:
//
//	typeSet.HaveNameMatching(".*Service")
func (ts *TypeSet) HaveNameMatching(pattern string) *TypeSet {
	return ts.NameMatch(pattern)
}

// AreInterfaces filters types that are interfaces
// It allows for filtering based on whether the type is an interface.
// Returns:
//   - *TypeSet: Returns the filtered TypeSet containing only interface types,
//     allowing for method chaining
//
// Example:
//
//	typeSet.AreInterfaces()
func (ts *TypeSet) AreInterfaces() *TypeSet {
	ts.currentPredicate = "AreInterfaces"

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		// Check if it's an interface using the IsInterface field
		if t.IsInterface {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}
