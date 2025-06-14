package goarchtest

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Types represents the entry point for architecture testing
type Types struct {
	pkgs    []*packages.Package
	typeSet *TypeSet
}

// TypeSet represents a collection of types that match certain criteria
type TypeSet struct {
	types             []*TypeInfo
	originalTypes     []*TypeInfo
	currentPredicate  string
	matchedPredicates []string
}

// TypeInfo contains comprehensive information about a Go type.
//
// This structure holds all the metadata needed for architectural analysis,
// including the type's name, location, dependencies, and structural characteristics.
//
// Fields:
//   - Name: The name of the type (e.g., "UserService")
//   - Package: The package name where the type is defined (e.g., "services")  
//   - FullPath: The full import path (e.g., "github.com/myorg/myapp/services")
//   - Imports: All import paths that this type's package depends on
//   - Interfaces: For interface types, the method names defined in the interface
//   - IsStruct: true if this type is a struct
//   - IsInterface: true if this type is an interface
//
// TypeInfo is used throughout GoArchTest's predicate system to make architectural
// decisions and validate constraints.
type TypeInfo struct {
	Name        string
	Package     string
	FullPath    string
	Imports     []string
	Interfaces  []string
	IsStruct    bool
	IsInterface bool
}

// InPath creates a new Types instance for packages in the specified directory path.
//
// This is the primary entry point for GoArchTest. It analyzes all Go packages
// found recursively in the given directory path and prepares them for architectural testing.
//
// Parameters:
//   - path: The directory path to analyze. Use "." for current directory or provide an absolute path.
//
// Returns:
//   - *Types: A Types instance containing all discovered types, ready for filtering and testing.
//
// Example:
//
//	// Analyze current project
//	types := goarchtest.InPath("./")
//	
//	// Analyze specific directory
//	types := goarchtest.InPath("/path/to/project")
//	
//	// Start testing architecture
//	result := types.That().
//	    ResideInNamespace("domain").
//	    ShouldNot().
//	    HaveDependencyOn("infrastructure").
//	    GetResult()
//
// The function uses Go's package loading mechanism to extract comprehensive
// type information including names, packages, imports, and structural details.
func InPath(path string) *Types {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports,
		Dir:  path,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load packages: %v\n", err)
		return &Types{
			pkgs:    []*packages.Package{},
			typeSet: &TypeSet{types: []*TypeInfo{}, originalTypes: []*TypeInfo{}},
		}
	}

	return &Types{
		pkgs:    pkgs,
		typeSet: extractTypesFromPackages(pkgs),
	}
}

// That starts a filter chain to select types
func (t *Types) That() *TypeSet {
	return t.typeSet.That()
}

// extractTypesFromPackages processes the packages to extract type information
func extractTypesFromPackages(pkgs []*packages.Package) *TypeSet {
	var types []*TypeInfo

	for _, pkg := range pkgs {
		// Skip packages with errors
		if len(pkg.Errors) > 0 {
			continue
		}

		imports := make([]string, 0)
		for importPath := range pkg.Imports {
			imports = append(imports, importPath)
		}

		// Get types from this package using syntax trees since we can't easily
		// map from types.Object to struct/interface information
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					typeInfo := &TypeInfo{
						Name:     typeSpec.Name.Name,
						Package:  pkg.Name,
						FullPath: pkg.PkgPath,
						Imports:  imports,
					}

					// Check if it's a struct
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						typeInfo.IsStruct = true
					}

					// Check if it's an interface
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						typeInfo.IsInterface = true
						// Collect method names from the interface
						if interfaceType.Methods != nil {
							for _, method := range interfaceType.Methods.List {
								if method.Names != nil {
									for _, name := range method.Names {
										typeInfo.Interfaces = append(typeInfo.Interfaces, name.Name)
									}
								}
							}
						}
					}

					types = append(types, typeInfo)
				}
			}
		}
	}

	return &TypeSet{
		types:         types,
		originalTypes: types,
	}
}

// That starts a filter chain
func (ts *TypeSet) That() *TypeSet {
	ts.currentPredicate = "That"
	return ts
}

// Result represents the outcome of architecture tests.
//
// It contains information about whether the architectural rule passed or failed,
// along with details about any types that violated the rule.
//
// Fields:
//   - IsSuccessful: true if the architectural test passed, false otherwise
//   - FailingTypes: slice of TypeInfo for types that didn't meet the criteria
//
// Example usage:
//
//	result := types.That().ResideInNamespace("domain").GetResult()
//	if !result.IsSuccessful {
//	    fmt.Printf("Found %d violations\n", len(result.FailingTypes))
//	    for _, failing := range result.FailingTypes {
//	        fmt.Printf("- %s in %s\n", failing.Name, failing.Package)
//	    }
//	}
type Result struct {
	IsSuccessful bool
	FailingTypes []*TypeInfo
}

// GetResult evaluates the predicates and returns the result
func (ts *TypeSet) GetResult() *Result {
	// If no predicates were applied, the test passes
	if len(ts.matchedPredicates) == 0 {
		return &Result{
			IsSuccessful: true,
		}
	}

	// Check if we have a negation in the predicates
	shouldNegate := false
	for _, pred := range ts.matchedPredicates {
		if pred == "Negate" {
			shouldNegate = true
			break
		}
	}

	// If we're negating, the result is successful if we have NO matching types
	if shouldNegate {
		return &Result{
			IsSuccessful: len(ts.types) == 0,
			FailingTypes: ts.types, // If we're negating, the failing types are the ones that matched
		}
	}

	// Otherwise, the result is successful if we have matching types
	return &Result{
		IsSuccessful: len(ts.types) > 0,
		FailingTypes: ts.getFailingTypes(),
	}
}

// getFailingTypes returns types that didn't match the predicates
func (ts *TypeSet) getFailingTypes() []*TypeInfo {
	var failingTypes []*TypeInfo

	// Compare original types with the filtered ones
	for _, origType := range ts.originalTypes {
		found := false
		for _, filteredType := range ts.types {
			if origType.FullPath == filteredType.FullPath && origType.Name == filteredType.Name {
				found = true
				break
			}
		}

		if !found {
			failingTypes = append(failingTypes, origType)
		}
	}

	return failingTypes
}

// GetAllTypes returns all types in the TypeSet
func (ts *TypeSet) GetAllTypes() []*TypeInfo {
	return ts.types
}

// GetFailureDetails returns a detailed error message about failing types
func (r *Result) GetFailureDetails() string {
	if r.IsSuccessful {
		return "No failures detected"
	}

	var details strings.Builder
	details.WriteString(fmt.Sprintf("Found %d failing type(s):\n", len(r.FailingTypes)))

	for i, failingType := range r.FailingTypes {
		details.WriteString(fmt.Sprintf("%d. %s in package %s\n", i+1, failingType.Name, failingType.Package))
	}

	return details.String()
}
