package examples

import (
	"fmt"
	"os"

	"github.com/solrac97gr/goarchtest"
)

func GenerateDependencyGraph() {
	// In a real application, you would use your project path
	projectPath := "/path/to/your/project"

	// Get all types in the project
	types := goarchtest.InPath(projectPath)

	// Create an error reporter
	reporter := goarchtest.NewErrorReporter(os.Stdout)

	// Get all types in the project
	allTypes := types.That().GetAllTypes()

	// Save the dependency graph to a dot file
	outputPath := "dependency_graph.dot"
	err := reporter.SaveDependencyGraph(allTypes, outputPath)
	if err != nil {
		fmt.Printf("Failed to save dependency graph: %v\n", err)
		return
	}

	fmt.Printf("Dependency graph saved to: %s\n", outputPath)
	fmt.Println("To generate a PNG, run: dot -Tpng dependency_graph.dot -o dependency_graph.png")
}
