package goarchtest

import (
	"fmt"
	"strings"
)

// ArchitecturePattern represents a predefined architectural pattern
type ArchitecturePattern struct {
	Name  string
	Rules []func(*Types) *Result
}

// Validate checks if the codebase adheres to the architectural pattern
func (ap *ArchitecturePattern) Validate(types *Types) []*ValidationResult {
	var results []*ValidationResult

	for i, rule := range ap.Rules {
		result := rule(types)
		validationResult := &ValidationResult{
			PatternName: ap.Name,
			RuleIndex:   i,
			IsSuccessful: result.IsSuccessful,
			FailingTypes: result.FailingTypes,
		}
		results = append(results, validationResult)
	}

	return results
}

// ValidationResult represents the result of validating an architectural pattern
type ValidationResult struct {
	PatternName  string
	RuleIndex    int
	IsSuccessful bool
	FailingTypes []*TypeInfo
}

// CleanArchitecture defines the Clean Architecture pattern (also known as Onion Architecture)
func CleanArchitecture(domainNamespace, applicationNamespace, infrastructureNamespace, presentationNamespace string) *ArchitecturePattern {
	return &ArchitecturePattern{
		Name: "Clean Architecture",
		Rules: []func(*Types) *Result{
			// Domain layer should not depend on any other layer
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(domainNamespace).
					ShouldNot().
					HaveDependencyOn(applicationNamespace).
					GetResult()
			},
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(domainNamespace).
					ShouldNot().
					HaveDependencyOn(infrastructureNamespace).
					GetResult()
			},
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(domainNamespace).
					ShouldNot().
					HaveDependencyOn(presentationNamespace).
					GetResult()
			},
			// Application layer should only depend on domain layer
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(applicationNamespace).
					ShouldNot().
					HaveDependencyOn(infrastructureNamespace).
					GetResult()
			},
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(applicationNamespace).
					ShouldNot().
					HaveDependencyOn(presentationNamespace).
					GetResult()
			},
			// Presentation layer should not depend on infrastructure layer
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(presentationNamespace).
					ShouldNot().
					HaveDependencyOn(infrastructureNamespace).
					GetResult()
			},
		},
	}
}

// HexagonalArchitecture defines the Hexagonal Architecture pattern (Ports and Adapters)
func HexagonalArchitecture(domainNamespace, portsNamespace, adaptersNamespace string) *ArchitecturePattern {
	return &ArchitecturePattern{
		Name: "Hexagonal Architecture",
		Rules: []func(*Types) *Result{
			// Domain should not depend on ports or adapters
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(domainNamespace).
					ShouldNot().
					HaveDependencyOn(portsNamespace).
					GetResult()
			},
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(domainNamespace).
					ShouldNot().
					HaveDependencyOn(adaptersNamespace).
					GetResult()
			},
			// Domain should be independent
			func(types *Types) *Result {
				domainTypes := types.That().ResideInNamespace(domainNamespace).types
				if len(domainTypes) == 0 {
					return &Result{
						IsSuccessful: false,
						FailingTypes: []*TypeInfo{
							{
								Name:    "Domain",
								Package: domainNamespace,
								FullPath: "No domain types found",
							},
						},
					}
				}
				return &Result{
					IsSuccessful: true,
				}
			},
			// Adapters should not be used directly by domain
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(adaptersNamespace).
					Should().
					ImplementInterface("Port").
					GetResult()
			},
		},
	}
}

// LayeredArchitecture defines a traditional layered architecture pattern
func LayeredArchitecture(layers ...string) *ArchitecturePattern {
	if len(layers) < 2 {
		panic("LayeredArchitecture requires at least 2 layers")
	}

	var rules []func(*Types) *Result

	// For each layer, ensure it doesn't depend on layers above it
	for i := 0; i < len(layers); i++ {
		currentLayer := layers[i]
		
		for j := i+1; j < len(layers); j++ {
			higherLayer := layers[j]
			
			// Create a closure to capture the layer values
			rule := func(current, higher string) func(*Types) *Result {
				return func(types *Types) *Result {
					return types.That().
						ResideInNamespace(current).
						ShouldNot().
						HaveDependencyOn(higher).
						GetResult()
				}
			}(currentLayer, higherLayer)
			
			rules = append(rules, rule)
		}
	}

	return &ArchitecturePattern{
		Name:  fmt.Sprintf("Layered Architecture (%s)", strings.Join(layers, " -> ")),
		Rules: rules,
	}
}

// MVCArchitecture defines the Model-View-Controller architecture pattern
func MVCArchitecture(modelNamespace, viewNamespace, controllerNamespace string) *ArchitecturePattern {
	return &ArchitecturePattern{
		Name: "MVC Architecture",
		Rules: []func(*Types) *Result{
			// Models should not depend on views or controllers
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(modelNamespace).
					ShouldNot().
					HaveDependencyOn(viewNamespace).
					GetResult()
			},
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(modelNamespace).
					ShouldNot().
					HaveDependencyOn(controllerNamespace).
					GetResult()
			},
			// Views should not depend on controllers
			func(types *Types) *Result {
				return types.That().
					ResideInNamespace(viewNamespace).
					ShouldNot().
					HaveDependencyOn(controllerNamespace).
					GetResult()
			},
		},
	}
}
