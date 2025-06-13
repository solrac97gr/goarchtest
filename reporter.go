package goarchtest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Reporter generates reports about architecture test results
type Reporter struct {
	Results []*Result
}

// NewReporter creates a new reporter instance
func NewReporter() *Reporter {
	return &Reporter{
		Results: []*Result{},
	}
}

// AddResult adds a test result to the reporter
func (r *Reporter) AddResult(result *Result) {
	r.Results = append(r.Results, result)
}

// GenerateTextReport generates a plain text report
func (r *Reporter) GenerateTextReport() string {
	var report strings.Builder
	
	report.WriteString("GoArchTest Report\n")
	report.WriteString("================\n\n")
	report.WriteString(fmt.Sprintf("Date: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	passCount := 0
	failCount := 0
	
	for i, result := range r.Results {
		if result.IsSuccessful {
			passCount++
			report.WriteString(fmt.Sprintf("Test #%d: PASS\n", i+1))
		} else {
			failCount++
			report.WriteString(fmt.Sprintf("Test #%d: FAIL\n", i+1))
			report.WriteString("Failing Types:\n")
			
			for _, failingType := range result.FailingTypes {
				report.WriteString(fmt.Sprintf("  - %s in package %s\n", failingType.Name, failingType.Package))
			}
			
			report.WriteString("\n")
		}
	}
	
	report.WriteString(fmt.Sprintf("\nSummary: %d passed, %d failed\n", passCount, failCount))
	
	return report.String()
}

// GenerateHTMLReport generates an HTML report
func (r *Reporter) GenerateHTMLReport() string {
	var report strings.Builder
	
	report.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>GoArchTest Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
        }
        h1 {
            color: #333;
        }
        .summary {
            margin: 20px 0;
            padding: 10px;
            background-color: #f5f5f5;
            border-radius: 4px;
        }
        .test {
            margin-bottom: 20px;
            padding: 10px;
            border-radius: 4px;
        }
        .pass {
            background-color: #dff0d8;
            border: 1px solid #d6e9c6;
        }
        .fail {
            background-color: #f2dede;
            border: 1px solid #ebccd1;
        }
        .test-title {
            font-weight: bold;
        }
        .failing-types {
            margin-top: 10px;
            margin-left: 20px;
        }
    </style>
</head>
<body>
    <h1>GoArchTest Report</h1>
    <p>Date: `)
	
	report.WriteString(time.Now().Format("2006-01-02 15:04:05"))
	report.WriteString(`</p>`)
	
	passCount := 0
	failCount := 0
	
	for i, result := range r.Results {
		if result.IsSuccessful {
			passCount++
			report.WriteString(fmt.Sprintf(`
    <div class="test pass">
        <div class="test-title">Test #%d: PASS</div>
    </div>`, i+1))
		} else {
			failCount++
			report.WriteString(fmt.Sprintf(`
    <div class="test fail">
        <div class="test-title">Test #%d: FAIL</div>
        <div class="failing-types">
            <strong>Failing Types:</strong>
            <ul>`, i+1))
			
			for _, failingType := range result.FailingTypes {
				report.WriteString(fmt.Sprintf(`
                <li>%s in package %s</li>`, failingType.Name, failingType.Package))
			}
			
			report.WriteString(`
            </ul>
        </div>
    </div>`)
		}
	}
	
	report.WriteString(fmt.Sprintf(`
    <div class="summary">
        <strong>Summary:</strong> %d passed, %d failed
    </div>
</body>
</html>`, passCount, failCount))
	
	return report.String()
}

// SaveReport saves a report to a file
func (r *Reporter) SaveReport(reportType string, outputPath string) error {
	var content string
	
	switch strings.ToLower(reportType) {
	case "text":
		content = r.GenerateTextReport()
	case "html":
		content = r.GenerateHTMLReport()
	default:
		return fmt.Errorf("unsupported report type: %s", reportType)
	}
	
	// Ensure the directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Write the report to file
	return os.WriteFile(outputPath, []byte(content), 0644)
}
