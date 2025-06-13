package goarchtest

// CustomPredicate represents a custom predicate function
type CustomPredicate func(*TypeInfo) bool

// WithCustomPredicate applies a custom predicate to the TypeSet
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
