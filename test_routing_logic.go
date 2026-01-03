package main

import (
	"fmt"
	"strings"
)

func main() {
	// Simulate the routing
	testPaths := []struct {
		path       string
		method     string
		expectedOK bool
	}{
		{path: "/exports", method: "GET", expectedOK: true},
		{path: "/exports/markdown", method: "POST", expectedOK: true},
		{path: "/exports/markdown", method: "GET", expectedOK: false},
		{path: "/series", method: "GET", expectedOK: true},
		{path: "/series/123", method: "GET", expectedOK: true},
	}

	fmt.Println("Testing routing logic:")
	fmt.Println("============================================================")

	for _, tc := range testPaths {
		fullPath := strings.TrimPrefix(tc.path, "api")
		fullPath = strings.Trim(fullPath, "/")

		parts := strings.Split(fullPath, "/")
		resource := parts[0]
		var id string
		if len(parts) > 1 {
			id = parts[1]
		}

		fmt.Printf("\nPath: %s %s\n", tc.method, tc.path)
		fmt.Printf("  resource: %q, id: %q\n", resource, id)

		// Determine handler call
		ok := false
		var handler string
		switch resource {
		case "exports":
			if id == "markdown" && tc.method == "POST" {
				handler = "handleExportsMarkdown()"
				ok = true
			} else if id == "" && tc.method == "GET" {
				handler = "handleExportsGet()"
				ok = true
			} else {
				handler = "405 Method Not Allowed"
				ok = false
			}
		case "series":
			handler = "handleSeries(id=" + id + ")"
			ok = true
		}

		status := "✓"
		if !ok {
			status = "✗"
		}
		fmt.Printf("  %s %s\n", status, handler)
		
		if ok != tc.expectedOK {
			fmt.Printf("  ERROR: expected ok=%v, got %v\n", tc.expectedOK, ok)
		}
	}
}
