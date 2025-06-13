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
func (g *GoArchTest) CheckRule(rule func(*Types) *Result) *Result {
	return rule(g.Types)
}
