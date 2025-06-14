package goarchtest

import (
	"fmt"
	"strings"
)

// Rule represents an architectural rule with a description and validation function
type Rule struct {
	Description string
	Validate    func(*Types) *Result
}

// ArchitecturePattern represents a predefined architectural pattern
type ArchitecturePattern struct {
	Name  string
	Rules []Rule
}

// Validate checks if the codebase adheres to the architectural pattern
func (ap *ArchitecturePattern) Validate(types *Types) []*ValidationResult {
	var results []*ValidationResult

	for i, rule := range ap.Rules {
		result := rule.Validate(types)
		validationResult := &ValidationResult{
			PatternName:     ap.Name,
			RuleIndex:       i,
			RuleDescription: rule.Description,
			IsSuccessful:    result.IsSuccessful,
			FailingTypes:    result.FailingTypes,
		}
		results = append(results, validationResult)
	}

	return results
}

// ValidationResult represents the result of validating an architectural pattern
type ValidationResult struct {
	PatternName     string
	RuleIndex       int
	RuleDescription string
	IsSuccessful    bool
	FailingTypes    []*TypeInfo
}

// CleanArchitecture defines the Clean Architecture pattern (also known as Onion Architecture)
func CleanArchitecture(domainNamespace, applicationNamespace, infrastructureNamespace, presentationNamespace string) *ArchitecturePattern {
	return &ArchitecturePattern{
		Name: "Clean Architecture",
		Rules: []Rule{
			// Domain layer should not depend on any other layer
			{
				Description: fmt.Sprintf("Domain layer (%s) should not depend on application layer (%s)", domainNamespace, applicationNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNamespace).
						ShouldNot().
						HaveDependencyOn(applicationNamespace).
						GetResult()
				},
			},
			{
				Description: fmt.Sprintf("Domain layer (%s) should not depend on infrastructure layer (%s)", domainNamespace, infrastructureNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNamespace).
						ShouldNot().
						HaveDependencyOn(infrastructureNamespace).
						GetResult()
				},
			},
			{
				Description: fmt.Sprintf("Domain layer (%s) should not depend on presentation layer (%s)", domainNamespace, presentationNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNamespace).
						ShouldNot().
						HaveDependencyOn(presentationNamespace).
						GetResult()
				},
			},
			// Application layer should only depend on domain layer
			{
				Description: fmt.Sprintf("Application layer (%s) should not depend on infrastructure layer (%s)", applicationNamespace, infrastructureNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(applicationNamespace).
						ShouldNot().
						HaveDependencyOn(infrastructureNamespace).
						GetResult()
				},
			},
			{
				Description: fmt.Sprintf("Application layer (%s) should not depend on presentation layer (%s)", applicationNamespace, presentationNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(applicationNamespace).
						ShouldNot().
						HaveDependencyOn(presentationNamespace).
						GetResult()
				},
			},
			// Presentation layer should not depend on infrastructure layer
			{
				Description: fmt.Sprintf("Presentation layer (%s) should not depend on infrastructure layer (%s)", presentationNamespace, infrastructureNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(presentationNamespace).
						ShouldNot().
						HaveDependencyOn(infrastructureNamespace).
						GetResult()
				},
			},
		},
	}
}

// HexagonalArchitecture defines the Hexagonal Architecture pattern (Ports and Adapters)
func HexagonalArchitecture(domainNamespace, portsNamespace, adaptersNamespace string) *ArchitecturePattern {
	return &ArchitecturePattern{
		Name: "Hexagonal Architecture",
		Rules: []Rule{
			// Domain should not depend on ports or adapters
			{
				Description: fmt.Sprintf("Domain layer (%s) should not depend on ports layer (%s)", domainNamespace, portsNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNamespace).
						ShouldNot().
						HaveDependencyOn(portsNamespace).
						GetResult()
				},
			},
			{
				Description: fmt.Sprintf("Domain layer (%s) should not depend on adapters layer (%s)", domainNamespace, adaptersNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNamespace).
						ShouldNot().
						HaveDependencyOn(adaptersNamespace).
						GetResult()
				},
			},
			// Domain should be independent
			{
				Description: fmt.Sprintf("Domain layer (%s) should exist", domainNamespace),
				Validate: func(types *Types) *Result {
					domainTypes := types.That().ResideInNamespace(domainNamespace).types
					if len(domainTypes) == 0 {
						return &Result{
							IsSuccessful: false,
							FailingTypes: []*TypeInfo{
								{
									Name:     "Domain",
									Package:  domainNamespace,
									FullPath: "No domain types found",
								},
							},
						}
					}
					return &Result{
						IsSuccessful: true,
					}
				},
			},
			// Adapters should not be used directly by domain
			{
				Description: fmt.Sprintf("Adapters (%s) should implement a Port interface", adaptersNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(adaptersNamespace).
						Should().
						ImplementInterface("Port").
						GetResult()
				},
			},
		},
	}
}

// LayeredArchitecture defines a traditional layered architecture pattern
func LayeredArchitecture(layers ...string) *ArchitecturePattern {
	if len(layers) < 2 {
		panic("LayeredArchitecture requires at least 2 layers")
	}

	var rules []Rule

	// For each layer, ensure it doesn't depend on layers above it
	for i := 0; i < len(layers); i++ {
		currentLayer := layers[i]

		for j := i + 1; j < len(layers); j++ {
			higherLayer := layers[j]

			// Create a rule with description and validation function
			rule := Rule{
				Description: fmt.Sprintf("Layer %s should not depend on higher layer %s", currentLayer, higherLayer),
				Validate: func(current, higher string) func(*Types) *Result {
					return func(types *Types) *Result {
						return types.That().
							ResideInNamespace(current).
							ShouldNot().
							HaveDependencyOn(higher).
							GetResult()
					}
				}(currentLayer, higherLayer),
			}

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
		Rules: []Rule{
			// Models should not depend on views or controllers
			{
				Description: fmt.Sprintf("Model layer (%s) should not depend on view layer (%s)", modelNamespace, viewNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(modelNamespace).
						ShouldNot().
						HaveDependencyOn(viewNamespace).
						GetResult()
				},
			},
			{
				Description: fmt.Sprintf("Model layer (%s) should not depend on controller layer (%s)", modelNamespace, controllerNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(modelNamespace).
						ShouldNot().
						HaveDependencyOn(controllerNamespace).
						GetResult()
				},
			},
			// Views should not depend on controllers
			{
				Description: fmt.Sprintf("View layer (%s) should not depend on controller layer (%s)", viewNamespace, controllerNamespace),
				Validate: func(types *Types) *Result {
					return types.That().
						ResideInNamespace(viewNamespace).
						ShouldNot().
						HaveDependencyOn(controllerNamespace).
						GetResult()
				},
			},
		},
	}
}

// DDDWithCleanArchitecture defines a Domain-Driven Design pattern with Clean Architecture within each bounded context
// This pattern enforces:
// 1. Bounded contexts are isolated from each other (no cross-domain dependencies)
// 2. Within each domain: Clean Architecture layers (domain -> application -> infrastructure)
// 3. Shared kernel can be used by all domains
// 4. pkg/ contains reusable utilities that can be used by any layer
func DDDWithCleanArchitecture(domains []string, sharedNamespace, pkgNamespace string) *ArchitecturePattern {
	var rules []Rule

	// Rule 1: Domain layers should not depend on application or infrastructure within the same domain
	for _, domain := range domains {
		domainNS := fmt.Sprintf("internal/%s/domain", domain)
		applicationNS := fmt.Sprintf("internal/%s/application", domain)
		infrastructureNS := fmt.Sprintf("internal/%s/infrastructure", domain)

		// Domain should not depend on application
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Domain layer (%s) should not depend on application layer (%s)", domainNS, applicationNS),
			Validate: func(domainNS, applicationNS string) func(*Types) *Result {
				return func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNS).
						ShouldNot().
						HaveDependencyOn(applicationNS).
						GetResult()
				}
			}(domainNS, applicationNS),
		})

		// Domain should not depend on infrastructure
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Domain layer (%s) should not depend on infrastructure layer (%s)", domainNS, infrastructureNS),
			Validate: func(domainNS, infrastructureNS string) func(*Types) *Result {
				return func(types *Types) *Result {
					return types.That().
						ResideInNamespace(domainNS).
						ShouldNot().
						HaveDependencyOn(infrastructureNS).
						GetResult()
				}
			}(domainNS, infrastructureNS),
		})

		// Application should not depend on infrastructure
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Application layer (%s) should not depend on infrastructure layer (%s)", applicationNS, infrastructureNS),
			Validate: func(applicationNS, infrastructureNS string) func(*Types) *Result {
				return func(types *Types) *Result {
					return types.That().
						ResideInNamespace(applicationNS).
						ShouldNot().
						HaveDependencyOn(infrastructureNS).
						GetResult()
				}
			}(applicationNS, infrastructureNS),
		})
	}

	// Rule 2: Cross-domain dependencies are not allowed (bounded context isolation)
	for i, domain1 := range domains {
		for j, domain2 := range domains {
			if i != j {
				domain1Prefix := fmt.Sprintf("internal/%s", domain1)
				domain2Prefix := fmt.Sprintf("internal/%s", domain2)

				rules = append(rules, Rule{
					Description: fmt.Sprintf("Domain %s should not depend on domain %s (bounded context isolation)", domain1, domain2),
					Validate: func(d1, d2 string) func(*Types) *Result {
						return func(types *Types) *Result {
							return types.That().
								ResideInNamespace(d1).
								ShouldNot().
								HaveDependencyOn(d2).
								GetResult()
						}
					}(domain1Prefix, domain2Prefix),
				})
			}
		}
	}

	// Rule 3: Only domain layers can depend on shared namespace (shared kernel)
	if sharedNamespace != "" {
		for _, domain := range domains {
			applicationNS := fmt.Sprintf("internal/%s/application", domain)
			infrastructureNS := fmt.Sprintf("internal/%s/infrastructure", domain)

			// Application should not depend on shared (only domain can)
			rules = append(rules, Rule{
				Description: fmt.Sprintf("Application layer (%s) should not depend on shared kernel (%s)", applicationNS, sharedNamespace),
				Validate: func(applicationNS, sharedNS string) func(*Types) *Result {
					return func(types *Types) *Result {
						return types.That().
							ResideInNamespace(applicationNS).
							ShouldNot().
							HaveDependencyOn(sharedNS).
							GetResult()
					}
				}(applicationNS, sharedNamespace),
			})

			// Infrastructure should not depend on shared (only domain can)
			rules = append(rules, Rule{
				Description: fmt.Sprintf("Infrastructure layer (%s) should not depend on shared kernel (%s)", infrastructureNS, sharedNamespace),
				Validate: func(infrastructureNS, sharedNS string) func(*Types) *Result {
					return func(types *Types) *Result {
						return types.That().
							ResideInNamespace(infrastructureNS).
							ShouldNot().
							HaveDependencyOn(sharedNS).
							GetResult()
					}
				}(infrastructureNS, sharedNamespace),
			})
		}
	}

	return &ArchitecturePattern{
		Name:  fmt.Sprintf("DDD with Clean Architecture (domains: %s)", strings.Join(domains, ", ")),
		Rules: rules,
	}
}

// CQRSArchitecture defines the Command Query Responsibility Segregation pattern
// This pattern enforces:
// 1. Commands (write operations) and Queries (read operations) are separated
// 2. Commands should not depend on queries
// 3. Queries should not depend on commands
// 4. Both can depend on shared domain models
// 5. Commands typically interact with write models/aggregates
// 6. Queries typically interact with read models/projections
func CQRSArchitecture(commandNamespace, queryNamespace, domainNamespace, writeModelNamespace, readModelNamespace string) *ArchitecturePattern {
	var rules []Rule

	// Rule 1: Commands should not depend on queries (separation of concerns)
	rules = append(rules, Rule{
		Description: fmt.Sprintf("Commands (%s) should not depend on queries (%s) - separation of concerns", commandNamespace, queryNamespace),
		Validate: func(types *Types) *Result {
			return types.That().
				ResideInNamespace(commandNamespace).
				ShouldNot().
				HaveDependencyOn(queryNamespace).
				GetResult()
		},
	})

	// Rule 2: Queries should not depend on commands (separation of concerns)
	rules = append(rules, Rule{
		Description: fmt.Sprintf("Queries (%s) should not depend on commands (%s) - separation of concerns", queryNamespace, commandNamespace),
		Validate: func(types *Types) *Result {
			return types.That().
				ResideInNamespace(queryNamespace).
				ShouldNot().
				HaveDependencyOn(commandNamespace).
				GetResult()
		},
	})

	// Rule 3: Write models should not depend on read models
	if writeModelNamespace != "" && readModelNamespace != "" {
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Write models (%s) should not depend on read models (%s)", writeModelNamespace, readModelNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(writeModelNamespace).
					ShouldNot().
					HaveDependencyOn(readModelNamespace).
					GetResult()
			},
		})

		// Rule 4: Read models should not depend on write models
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Read models (%s) should not depend on write models (%s)", readModelNamespace, writeModelNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(readModelNamespace).
					ShouldNot().
					HaveDependencyOn(writeModelNamespace).
					GetResult()
			},
		})

		// Rule 5: Commands should primarily use write models
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Commands (%s) should not depend on read models (%s)", commandNamespace, readModelNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(commandNamespace).
					ShouldNot().
					HaveDependencyOn(readModelNamespace).
					GetResult()
			},
		})

		// Rule 6: Queries should primarily use read models
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Queries (%s) should not depend on write models (%s)", queryNamespace, writeModelNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(queryNamespace).
					ShouldNot().
					HaveDependencyOn(writeModelNamespace).
					GetResult()
			},
		})
	}

	// Rule 7: Both commands and queries can depend on shared domain (if provided)
	// This is allowed, so no restriction rule needed

	return &ArchitecturePattern{
		Name:  "CQRS Architecture",
		Rules: rules,
	}
}

// EventSourcedCQRSArchitecture defines CQRS with Event Sourcing pattern
// This pattern enforces:
// 1. All CQRS rules
// 2. Commands should depend on events (to produce them)
// 3. Event store is the source of truth for commands
// 4. Queries should not depend on event store directly (use projections)
// 5. Projections should depend on events (to build read models)
func EventSourcedCQRSArchitecture(commandNamespace, queryNamespace, eventNamespace, eventStoreNamespace, projectionNamespace, domainNamespace string) *ArchitecturePattern {
	var rules []Rule

	// Include basic CQRS rules
	cqrsPattern := CQRSArchitecture(commandNamespace, queryNamespace, domainNamespace, "", "")
	rules = append(rules, cqrsPattern.Rules...)

	// Rule 1: Commands should have dependency on events namespace (to produce them)
	if eventNamespace != "" {
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Commands (%s) should depend on events (%s) to produce them", commandNamespace, eventNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(commandNamespace).
					Should().
					HaveDependencyOn(eventNamespace).
					GetResult()
			},
		})
	}

	// Rule 2: Commands should interact with event store
	if eventStoreNamespace != "" {
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Commands (%s) should depend on event store (%s)", commandNamespace, eventStoreNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(commandNamespace).
					Should().
					HaveDependencyOn(eventStoreNamespace).
					GetResult()
			},
		})
	}

	// Rule 3: Queries should not depend on event store directly (use projections instead)
	if eventStoreNamespace != "" {
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Queries (%s) should not depend on event store (%s) directly", queryNamespace, eventStoreNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(queryNamespace).
					ShouldNot().
					HaveDependencyOn(eventStoreNamespace).
					GetResult()
			},
		})
	}

	// Rule 4: Projections should depend on events (to build read models)
	if projectionNamespace != "" && eventNamespace != "" {
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Projections (%s) should depend on events (%s) to build read models", projectionNamespace, eventNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(projectionNamespace).
					Should().
					HaveDependencyOn(eventNamespace).
					GetResult()
			},
		})
	}

	// Rule 5: Queries should depend on projections (not directly on events)
	if projectionNamespace != "" {
		rules = append(rules, Rule{
			Description: fmt.Sprintf("Queries (%s) should depend on projections (%s) not directly on events", queryNamespace, projectionNamespace),
			Validate: func(types *Types) *Result {
				return types.That().
					ResideInNamespace(queryNamespace).
					Should().
					HaveDependencyOn(projectionNamespace).
					GetResult()
			},
		})
	}

	return &ArchitecturePattern{
		Name:  "Event Sourced CQRS Architecture",
		Rules: rules,
	}
}
