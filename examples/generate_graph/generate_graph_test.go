package generate_graph_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestGenerateDependencyGraph(t *testing.T) {
	// Get the path to the sample project
	sampleProjectPath, err := filepath.Abs("../sample_project")
	if err != nil {
		t.Fatalf("Failed to get sample project path: %v", err)
	}

	// Create a new GoArchTest instance
	types := goarchtest.InPath(sampleProjectPath)

	// Create an error reporter that writes to stderr
	reporter := goarchtest.NewErrorReporter(os.Stderr)

	// Get all types in the sample project
	allTypes := types.That().GetAllTypes()

	// Create the output directory if it doesn't exist
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Save the dependency graph to a dot file
	dotFilePath := filepath.Join(currentDir, "dependency_graph.dot")
	err = reporter.SaveDependencyGraph(allTypes, dotFilePath)
	if err != nil {
		t.Fatalf("Failed to save dependency graph: %v", err)
	}

	t.Logf("Dependency graph saved to: %s", dotFilePath)

	// Check if Graphviz dot is installed
	_, err = exec.LookPath("dot")
	if err != nil {
		t.Logf("Graphviz not found. To visualize the graph, install Graphviz and run:")
		t.Logf("  dot -Tpng %s -o dependency_graph.png", dotFilePath)
	} else {
		// Generate PNG image from dot file
		pngFilePath := filepath.Join(currentDir, "dependency_graph.png")
		cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", pngFilePath)
		if err := cmd.Run(); err != nil {
			t.Logf("Failed to generate PNG: %v", err)
			t.Logf("To generate manually, run: dot -Tpng %s -o dependency_graph.png", dotFilePath)
		} else {
			t.Logf("PNG image generated at: %s", pngFilePath)

			// Try to open the image
			var openCmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin":
				openCmd = exec.Command("open", pngFilePath)
			case "windows":
				openCmd = exec.Command("cmd", "/c", "start", pngFilePath)
			case "linux":
				openCmd = exec.Command("xdg-open", pngFilePath)
			default:
				t.Logf("Unsupported OS for automatic image opening")
			}

			if openCmd != nil {
				if err := openCmd.Start(); err != nil {
					t.Logf("Failed to open the image: %v", err)
				}
			}
		}
	}

	// Display instructions for interpreting the graph
	t.Log("\nHow to interpret the dependency graph:")
	t.Log("- Each box represents a package in your project")
	t.Log("- Arrows show dependencies between packages")
	t.Log("- The direction of the arrow indicates which package depends on another")
	t.Log("- For clean architecture, domain should have no outgoing arrows")
}

// Example provides a simple example of how to generate a dependency graph
func Example() {
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
