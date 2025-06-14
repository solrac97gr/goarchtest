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

// TypeInfo contains information about a Go type
type TypeInfo struct {
	Name        string
	Package     string
	FullPath    string
	Imports     []string
	Interfaces  []string
	IsStruct    bool
	IsInterface bool
}

// InPath creates a new Types instance for packages in the specified directory path
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

// Result represents the outcome of architecture tests
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
