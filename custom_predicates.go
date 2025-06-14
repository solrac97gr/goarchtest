package goarchtest

// CustomPredicate represents a custom predicate function
type CustomPredicate func(*TypeInfo) bool

// WithCustomPredicate applies a custom predicate to the TypeSet
// WithCustomPredicate applies a custom predicate function to filter the TypeSet.
// It allows for flexible filtering of types based on custom logic defined by the caller.
//
// Parameters:
//   - name: A string identifier for the predicate being applied
//   - predicate: A CustomPredicate function that takes a *TypeInfo and returns bool
//
// The predicate function is called for each type in the set. Only types that return
// true when passed to the predicate are kept in the filtered set.
//
// Returns *TypeSet to allow for method chaining.
//
// Example:
//
//	ts.WithCustomPredicate("hasFields", func(t *TypeInfo) bool {
//	    return len(t.Fields) > 0
//	})
func (ts *TypeSet) WithCustomPredicate(name string, predicate CustomPredicate) *TypeSet {
	ts.currentPredicate = name

	var filteredTypes []*TypeInfo
	for _, t := range ts.types {
		if predicate(t) {
			filteredTypes = append(filteredTypes, t)
		}
	}

	ts.types = filteredTypes
	ts.matchedPredicates = append(ts.matchedPredicates, ts.currentPredicate)
	return ts
}
