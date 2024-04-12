package cmd

import (
	"testing"
)

func TestIsBallerinaProject(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"../testData/BalFile/main.bal", false},
		{"../testData/BalProject", true},
	}

	for _, tc := range testCases {
		actual := isBallerinaProject(tc.path)
		if actual != tc.expected {
			t.Errorf("Expected isBallerinaProject(%s) to be %v, but got %v", tc.path, tc.expected, actual)
		}
	}
}

func TestGetExpectedOuput(t *testing.T) {
	testCases := []struct {
		path     string
		expected string
	}{
		{"../testData/BalFile/main.bal", "main.jar"},
		{"../testData/BalProject", "BalProject.jar"},
	}

	for _, tc := range testCases {
		actual := getExpectedOutput(tc.path)
		if actual != tc.expected {
			t.Errorf("Expected getExpectedOutput(%s) to be %s, but got %s", tc.path, tc.expected, actual)
		}
	}
}
