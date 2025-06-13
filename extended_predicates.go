package goarchtest

import (
	"regexp"
	"strings"
)

// NameMatch filters types based on a regex pattern match on their names
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

// HaveNameEndingWith filters types whose names end with the specified suffix
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

// DoNotHaveDependencyOn filters types that do not have a dependency on the specified package
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
