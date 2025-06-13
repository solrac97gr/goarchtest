package generate_graph

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/solrac97gr/goarchtest"
)

func GenerateDependencyGraph() {
	// Parse command-line flags
	projectPath := flag.String("path", ".", "Path to the project directory")
	outputFile := flag.String("output", "dependency_graph.dot", "Output path for the DOT file")
	generatePNG := flag.Bool("png", false, "Generate PNG from DOT file")
	openImage := flag.Bool("open", false, "Open the generated PNG image")
	flag.Parse()

	// Convert to absolute path
	absPath, err := filepath.Abs(*projectPath)
	if err != nil {
		fmt.Printf("Error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Analyzing project at: %s\n", absPath)

	// Create a new GoArchTest instance
	types := goarchtest.InPath(absPath)

	// Create an error reporter
	reporter := goarchtest.NewErrorReporter(os.Stdout)

	// Get all types in the project
	allTypes := types.That().GetAllTypes()

	fmt.Printf("Found %d types in the project\n", len(allTypes))

	// Save the dependency graph to a dot file
	err = reporter.SaveDependencyGraph(allTypes, *outputFile)
	if err != nil {
		fmt.Printf("Failed to save dependency graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Dependency graph saved to: %s\n", *outputFile)

	// Generate PNG if requested
	if *generatePNG {
		pngFilePath := filepath.Join(filepath.Dir(*outputFile),
			fmt.Sprintf("%s.png", strings.TrimSuffix(filepath.Base(*outputFile), filepath.Ext(*outputFile))))

		// Check if Graphviz dot is installed
		_, err = exec.LookPath("dot")
		if err != nil {
			fmt.Println("Graphviz not found. Please install Graphviz to generate PNG images.")
			fmt.Printf("Then run: dot -Tpng %s -o %s\n", *outputFile, pngFilePath)
			os.Exit(1)
		}

		// Generate PNG image from dot file
		cmd := exec.Command("dot", "-Tpng", *outputFile, "-o", pngFilePath)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to generate PNG: %v\n", err)
			fmt.Printf("To generate manually, run: dot -Tpng %s -o %s\n", *outputFile, pngFilePath)
			os.Exit(1)
		}

		fmt.Printf("PNG image generated at: %s\n", pngFilePath)

		// Open the image if requested
		if *openImage {
			var openCmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin":
				openCmd = exec.Command("open", pngFilePath)
			case "windows":
				openCmd = exec.Command("cmd", "/c", "start", pngFilePath)
			case "linux":
				openCmd = exec.Command("xdg-open", pngFilePath)
			default:
				fmt.Println("Unsupported OS for automatic image opening")
				openCmd = nil
			}

			if openCmd != nil {
				if err := openCmd.Start(); err != nil {
					fmt.Printf("Failed to open the image: %v\n", err)
				}
			}
		}
	} else {
		fmt.Printf("To generate a PNG, install Graphviz and run: dot -Tpng %s -o dependency_graph.png\n", *outputFile)
	}

	// Display instructions for interpreting the graph
	fmt.Println("\nHow to interpret the dependency graph:")
	fmt.Println("- Each box represents a package in your project")
	fmt.Println("- Arrows show dependencies between packages")
	fmt.Println("- The direction of the arrow indicates which package depends on another")
	fmt.Println("- For clean architecture, domain should have no outgoing arrows")
}
