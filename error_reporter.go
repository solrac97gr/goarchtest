package goarchtest

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ErrorReporter handles reporting of architecture test errors
type ErrorReporter struct {
	writer io.Writer
}

// NewErrorReporter creates a new ErrorReporter with the specified writer
func NewErrorReporter(writer io.Writer) *ErrorReporter {
	if writer == nil {
		writer = os.Stdout
	}
	return &ErrorReporter{
		writer: writer,
	}
}

// ReportError reports an error from an architecture test
func (er *ErrorReporter) ReportError(result *Result, description string) {
	if result.IsSuccessful {
		return
	}

	fmt.Fprintf(er.writer, "Architecture Test Failed: %s\n", description)

	if len(result.FailingTypes) > 0 {
		fmt.Fprintln(er.writer, "Failing Types:")

		for _, failingType := range result.FailingTypes {
			fmt.Fprintf(er.writer, "  - %s in package %s\n", failingType.Name, failingType.Package)
		}
	}

	fmt.Fprintln(er.writer)
}

// ReportPatternValidation reports the results of validating an architectural pattern
func (er *ErrorReporter) ReportPatternValidation(results []*ValidationResult) {
	if len(results) == 0 {
		return
	}

	patternName := results[0].PatternName
	fmt.Fprintf(er.writer, "Validating %s Pattern\n", patternName)
	fmt.Fprintf(er.writer, "%s\n", strings.Repeat("=", len(patternName)+18))

	passCount := 0
	failCount := 0

	for i, result := range results {
		if result.IsSuccessful {
			passCount++
			fmt.Fprintf(er.writer, "Rule #%d: PASS\n", i+1)
		} else {
			failCount++
			fmt.Fprintf(er.writer, "Rule #%d: FAIL\n", i+1)

			if len(result.FailingTypes) > 0 {
				fmt.Fprintln(er.writer, "Failing Types:")

				for _, failingType := range result.FailingTypes {
					fmt.Fprintf(er.writer, "  - %s in package %s\n", failingType.Name, failingType.Package)
				}
			}

			fmt.Fprintln(er.writer)
		}
	}

	fmt.Fprintf(er.writer, "\nSummary: %d/%d rules passed\n", passCount, len(results))

	if failCount == 0 {
		fmt.Fprintf(er.writer, "The codebase adheres to the %s pattern.\n", patternName)
	} else {
		fmt.Fprintf(er.writer, "The codebase does NOT fully adhere to the %s pattern.\n", patternName)
	}

	fmt.Fprintln(er.writer)
}

// GenerateDependencyGraph generates a dot graph representing dependencies
// This can be used with Graphviz to visualize the dependencies
func (er *ErrorReporter) GenerateDependencyGraph(types []*TypeInfo) string {
	var graph strings.Builder

	graph.WriteString("digraph ArchitectureDependencies {\n")
	graph.WriteString("  rankdir=TB;\n")
	graph.WriteString("  node [shape=box, style=filled, fillcolor=lightblue];\n")

	// Handle nil or empty types
	if len(types) == 0 {
		graph.WriteString("  note [label=\"No types provided or empty type set\", shape=note, fillcolor=lightyellow];\n")
		graph.WriteString("}\n")
		return graph.String()
	}

	// Map of package names to node IDs
	packageNodes := make(map[string]string)
	nodeID := 0

	// Create nodes for each package
	for _, t := range types {
		if _, exists := packageNodes[t.Package]; !exists {
			packageNodes[t.Package] = fmt.Sprintf("node%d", nodeID)
			graph.WriteString(fmt.Sprintf("  %s [label=\"%s\"];\n", packageNodes[t.Package], t.Package))
			nodeID++
		}
	}

	// Create edges for dependencies
	for _, t := range types {
		srcNode := packageNodes[t.Package]

		// Add edges for each import
		for _, imp := range t.Imports {
			// Find if any of our packages match this import
			for pkg, dstNode := range packageNodes {
				if strings.Contains(imp, pkg) && srcNode != dstNode {
					graph.WriteString(fmt.Sprintf("  %s -> %s;\n", srcNode, dstNode))
					break
				}
			}
		}
	}

	graph.WriteString("}\n")

	return graph.String()
}

// SaveDependencyGraph saves a dependency graph to a file
func (er *ErrorReporter) SaveDependencyGraph(types []*TypeInfo, outputPath string) error {
	// If types is nil, we'll provide an informative message
	if len(types) == 0 {
		// Create a simple placeholder graph
		placeholderGraph := `digraph ArchitectureDependencies {
  rankdir=TB;
  node [shape=box, style=filled, fillcolor=lightblue];
  note [label="No types provided or empty type set", shape=note, fillcolor=lightyellow];
}
`
		return os.WriteFile(outputPath, []byte(placeholderGraph), 0644)
	}

	graph := er.GenerateDependencyGraph(types)
	return os.WriteFile(outputPath, []byte(graph), 0644)
}
