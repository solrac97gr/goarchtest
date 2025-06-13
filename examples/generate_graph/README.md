# Dependency Graph Generation Example

This example demonstrates how to generate dependency graphs for your Go projects using GoArchTest.

## Running the Test Example

```bash
go test -v
```

This will analyze the sample project and generate a dependency graph in DOT format. If you have Graphviz installed, it will also generate a PNG image and try to open it.

## Using the Command Line Tool

You can build and use the command line tool in the `cmd` directory:

```bash
# Build the tool
cd cmd
go build -o graph-generator

# Run the tool
./graph-generator
```

## How to Interpret the Graph

The dependency graph shows the relationships between packages in your project:

- Each box represents a package in your project
- Arrows show dependencies between packages
- The direction of the arrow indicates which package depends on another
- For clean architecture, domain should have no outgoing arrows

## Integrating in Your Own Code

Here's a minimal example of how to generate a dependency graph in your own code:

```go
package main

import (
	"fmt"
	"os"

	"github.com/solrac97gr/goarchtest"
)

func main() {
	// Path to your project
	projectPath := "."
	
	// Create a new GoArchTest instance
	types := goarchtest.InPath(projectPath)
	
	// Create an error reporter
	reporter := goarchtest.NewErrorReporter(os.Stdout)
	
	// Get all types in the project
	allTypes := types.That().GetAllTypes()
	
	// Save the dependency graph
	err := reporter.SaveDependencyGraph(allTypes, "dependency_graph.dot")
	if err != nil {
		fmt.Printf("Failed to save dependency graph: %v\n", err)
		return
	}
	
	fmt.Println("Dependency graph saved to: dependency_graph.dot")
	fmt.Println("To generate a PNG, run: dot -Tpng dependency_graph.dot -o dependency_graph.png")
}
```
