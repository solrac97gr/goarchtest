package goarchtest_test

import (
	"strings"
	"testing"

	"github.com/solrac97gr/goarchtest"
)

func TestGetFailureDetails(t *testing.T) {
	// Create a failing result
	failingResult := &goarchtest.Result{
		IsSuccessful: false,
		FailingTypes: []*goarchtest.TypeInfo{
			{
				Name:    "TestType1",
				Package: "test/package1",
			},
			{
				Name:    "TestType2",
				Package: "test/package2",
			},
		},
	}

	// Get the failure details
	details := failingResult.GetFailureDetails()

	// Check that the details contain the expected information
	if !strings.Contains(details, "Found 2 failing type(s)") {
		t.Errorf("Expected details to contain count of failing types, got: %s", details)
	}

	if !strings.Contains(details, "TestType1 in package test/package1") {
		t.Errorf("Expected details to contain first failing type, got: %s", details)
	}

	if !strings.Contains(details, "TestType2 in package test/package2") {
		t.Errorf("Expected details to contain second failing type, got: %s", details)
	}

	// Test with a successful result
	successResult := &goarchtest.Result{
		IsSuccessful: true,
		FailingTypes: nil,
	}

	successDetails := successResult.GetFailureDetails()
	if successDetails != "No failures detected" {
		t.Errorf("Expected 'No failures detected' for successful result, got: %s", successDetails)
	}
}
