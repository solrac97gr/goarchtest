package goarchtest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Types represents the entry point for architecture testing
type Types struct {
	packages map[string]*ast.Package
	typeSet  *TypeSet
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
	Name       string
	Package    string
	FullPath   string
	Imports    []string
	Interfaces []string
	IsStruct   bool
}

// InPath creates a new Types instance for packages in the specified directory path
func InPath(path string) *Types {
	packages := make(map[string]*ast.Package)
	fset := token.NewFileSet()

	// Walk through all Go files in the specified path
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Parse the directory
			pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// Add all packages to our map
			for name, pkg := range pkgs {
				packages[name] = pkg
			}
		}
		return nil
	})

	// Create a Types instance with the found packages
	return &Types{
		packages: packages,
		typeSet:  extractTypesFromPackages(packages),
	}
}

// That starts a filter chain to select types
func (t *Types) That() *TypeSet {
	return t.typeSet.That()
}

// extractTypesFromPackages processes the AST to extract type information
func extractTypesFromPackages(packages map[string]*ast.Package) *TypeSet {
	var types []*TypeInfo

	for pkgName, pkg := range packages {
		for fileName, file := range pkg.Files {
			// Extract imports
			imports := make([]string, 0)
			for _, imp := range file.Imports {
				if imp.Path != nil {
					// Remove quotes from import path
					impPath := strings.Trim(imp.Path.Value, "\"")
					imports = append(imports, impPath)
				}
			}

			// Find types in the file
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							typeInfo := &TypeInfo{
								Name:     typeSpec.Name.Name,
								Package:  pkgName,
								FullPath: fileName,
								Imports:  imports,
							}

							// Check if it's a struct
							if _, ok := typeSpec.Type.(*ast.StructType); ok {
								typeInfo.IsStruct = true
							}

							// Check if it implements interfaces
							if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
								for _, method := range interfaceType.Methods.List {
									for _, name := range method.Names {
										typeInfo.Interfaces = append(typeInfo.Interfaces, name.Name)
									}
								}
							}

							types = append(types, typeInfo)
						}
					}
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

	// The result is successful if we have matching types
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
