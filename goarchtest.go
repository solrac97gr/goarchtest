package goarchtest

// GoArchTest is the main entry point for the architecture testing library
type GoArchTest struct {
	Types *Types
}

// New creates a new instance of GoArchTest for the specified path
func New(path string) *GoArchTest {
	return &GoArchTest{
		Types: InPath(path),
	}
}

// CheckRule executes a predefined architectural rule
// It takes a rule function that operates on the Types and returns a Result.
// This function is used to apply various architectural checks on the types defined in the Go project.
// Parameters:
//   - rule: A function that takes a pointer to Types and returns a Result
//
// Returns:
//   - *Result: Returns the result of applying the rule, which includes whether the rule passed or failed
//
// Example:
//
//	result := goarchtest.CheckRule(goarchtest.ResideInNamespace("github.com/myorg/mypackage"))
func (g *GoArchTest) CheckRule(rule func(*Types) *Result) *Result {
	return rule(g.Types)
}
