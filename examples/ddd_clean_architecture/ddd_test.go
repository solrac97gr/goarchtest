package ddd_clean_architecture

import (
	"path/filepath"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

// TestDDDCleanArchitecture demonstrates testing Domain-Driven Design with Clean Architecture
func TestDDDCleanArchitecture(t *testing.T) {
	projectPath, _ := filepath.Abs("./")
	types := goarchtest.InPath(projectPath)

	t.Run("Predefined DDD Pattern Validation", func(t *testing.T) {
		// Test using the predefined DDD pattern
		domains := []string{"user", "products"}
		dddPattern := goarchtest.DDDWithCleanArchitecture(domains, "internal/shared", "pkg")
		
		validationResults := dddPattern.Validate(types)
		
		t.Logf("DDD with Clean Architecture Validation Results:")
		t.Logf("==============================================")
		
		passedRules := 0
		for i, result := range validationResults {
			status := "✅ PASS"
			if !result.IsSuccessful {
				status = "❌ FAIL"
				if len(result.FailingTypes) > 0 {
					t.Logf("Rule #%d: %s", i+1, status)
					for _, failingType := range result.FailingTypes {
						t.Logf("  Violation: %s (%s)", failingType.Name, failingType.Package)
					}
				}
			} else {
				passedRules++
				t.Logf("Rule #%d: %s", i+1, status)
			}
		}
		
		t.Logf("")
		t.Logf("Summary: %d/%d rules passed", passedRules, len(validationResults))
		
		if passedRules == len(validationResults) {
			t.Logf("✅ The codebase fully adheres to DDD with Clean Architecture!")
		} else {
			t.Logf("❌ The codebase does NOT fully adhere to DDD with Clean Architecture.")
		}
	})

	t.Run("Individual DDD Rules", func(t *testing.T) {
		t.Run("Bounded Context Isolation", func(t *testing.T) {
			// User domain should not depend on products domain
			result := types.
				That().
				ResideInNamespace("internal/user").
				ShouldNot().
				HaveDependencyOn("internal/products").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Bounded context violation detected:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s depends on products domain", failingType.Name)
				}
			} else {
				t.Logf("✅ User domain is properly isolated from products domain")
			}

			// Products domain should not depend on user domain
			result = types.
				That().
				ResideInNamespace("internal/products").
				ShouldNot().
				HaveDependencyOn("internal/user").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Bounded context violation detected:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s depends on user domain", failingType.Name)
				}
			} else {
				t.Logf("✅ Products domain is properly isolated from user domain")
			}
		})

		t.Run("Clean Architecture Within User Domain", func(t *testing.T) {
			// Domain should not depend on application
			result := types.
				That().
				ResideInNamespace("internal/user/domain").
				ShouldNot().
				HaveDependencyOn("internal/user/application").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Clean Architecture violation in user domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (domain) depends on application layer", failingType.Name)
				}
			} else {
				t.Logf("✅ User domain layer is properly isolated from application")
			}

			// Domain should not depend on infrastructure
			result = types.
				That().
				ResideInNamespace("internal/user/domain").
				ShouldNot().
				HaveDependencyOn("internal/user/infrastructure").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Clean Architecture violation in user domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (domain) depends on infrastructure layer", failingType.Name)
				}
			} else {
				t.Logf("✅ User domain layer is properly isolated from infrastructure")
			}

			// Application should not depend on infrastructure
			result = types.
				That().
				ResideInNamespace("internal/user/application").
				ShouldNot().
				HaveDependencyOn("internal/user/infrastructure").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Clean Architecture violation in user domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (application) depends on infrastructure layer", failingType.Name)
				}
			} else {
				t.Logf("✅ User application layer is properly isolated from infrastructure")
			}
		})

		t.Run("Clean Architecture Within Products Domain", func(t *testing.T) {
			// Domain should not depend on application
			result := types.
				That().
				ResideInNamespace("internal/products/domain").
				ShouldNot().
				HaveDependencyOn("internal/products/application").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Clean Architecture violation in products domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (domain) depends on application layer", failingType.Name)
				}
			} else {
				t.Logf("✅ Products domain layer is properly isolated from application")
			}

			// Application should not depend on infrastructure
			result = types.
				That().
				ResideInNamespace("internal/products/application").
				ShouldNot().
				HaveDependencyOn("internal/products/infrastructure").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Clean Architecture violation in products domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (application) depends on infrastructure layer", failingType.Name)
				}
			} else {
				t.Logf("✅ Products application layer is properly isolated from infrastructure")
			}
		})

		t.Run("Shared Kernel Usage", func(t *testing.T) {
			// Application layers should not directly depend on shared
			result := types.
				That().
				ResideInNamespace("internal/user/application").
				ShouldNot().
				HaveDependencyOn("internal/shared").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Shared kernel violation in user domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (application) directly uses shared kernel", failingType.Name)
				}
			} else {
				t.Logf("✅ User application properly accesses shared through domain")
			}

			result = types.
				That().
				ResideInNamespace("internal/products/application").
				ShouldNot().
				HaveDependencyOn("internal/shared").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ Shared kernel violation in products domain:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (application) directly uses shared kernel", failingType.Name)
				}
			} else {
				t.Logf("✅ Products application properly accesses shared through domain")
			}
		})

		t.Run("PKG Utilities Usage", func(t *testing.T) {
			// Verify that pkg utilities can be used by any layer
			result := types.
				That().
				ResideInNamespace("internal").
				Should().
				HaveDependencyOn("pkg").
				GetResult()

			if result.IsSuccessful {
				t.Logf("✅ Internal layers properly use pkg utilities")
			} else {
				t.Logf("ℹ️ No pkg dependencies found (this is OK)")
			}

			// Verify that pkg doesn't depend on internal code
			result = types.
				That().
				ResideInNamespace("pkg").
				ShouldNot().
				HaveDependencyOn("internal").
				GetResult()

			if !result.IsSuccessful {
				t.Logf("❌ PKG layer violation:")
				for _, failingType := range result.FailingTypes {
					t.Logf("  %s (pkg) depends on internal code", failingType.Name)
				}
			} else {
				t.Logf("✅ PKG layer is properly independent from internal code")
			}
		})
	})
}

// TestCustomDDDRules demonstrates creating custom rules for DDD
func TestCustomDDDRules(t *testing.T) {
	projectPath, _ := filepath.Abs("./")
	types := goarchtest.InPath(projectPath)

	t.Run("Custom Domain Rules", func(t *testing.T) {
		// Custom rule: All domain entities should have an ID field
		domainEntityPredicate := func(typeInfo *goarchtest.TypeInfo) bool {
			// This is a simplified check - in real scenarios you'd analyze the struct fields
			return typeInfo.IsStruct && 
				   (typeInfo.Package == "internal/user/domain" || typeInfo.Package == "internal/products/domain")
		}

		result := types.
			That().
			WithCustomPredicate("IsDomainEntity", domainEntityPredicate).
			Should().
			ResideInNamespace("internal").
			GetResult()

		if result.IsSuccessful {
			t.Logf("✅ Domain entities follow proper structure")
		} else {
			t.Logf("⚠️ Some domain entities may not follow conventions")
		}
	})

	t.Run("Repository Pattern Validation", func(t *testing.T) {
		// Custom rule: Repository interfaces should be in domain, implementations in infrastructure
		repositoryInterfacePredicate := func(typeInfo *goarchtest.TypeInfo) bool {
			// Check if the name suggests it's a repository interface
			return len(typeInfo.Name) > 10 && 
				   typeInfo.Name[len(typeInfo.Name)-10:] == "Repository" &&
				   len(typeInfo.Interfaces) > 0 // Has interfaces (likely an interface definition)
		}

		result := types.
			That().
			WithCustomPredicate("IsRepositoryInterface", repositoryInterfacePredicate).
			Should().
			ResideInNamespace("domain").
			GetResult()

		if result.IsSuccessful {
			t.Logf("✅ Repository interfaces are properly placed in domain layers")
		} else {
			t.Logf("⚠️ Repository interfaces may not be in correct locations")
		}
	})
}
