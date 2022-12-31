package main

import (
	"fmt"
	"strconv"
	"strings"
)

// isGreater returns true if v1 is greater than v2 according to semantic versioning rules.
func isGreater(v1, v2 string) bool {
	// parse the version numbers
	major1, minor1, patch1, err := parseVersion(v1)
	if err != nil {
		return false
	}
	major2, minor2, patch2, err := parseVersion(v2)
	if err != nil {
		// handle error
		return false
	}

	// compare the versions
	if major1 > major2 {
		return true
	}
	if major1 < major2 {
		return false
	}
	if minor1 > minor2 {
		return true
	}
	if minor1 < minor2 {
		return false
	}
	if patch1 > patch2 {
		return true
	}
	return false
}

// parseVersion parses a semantic version number and returns its components as integers.
func parseVersion(v string) (int, int, int, error) {
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("invalid version number: %s", v)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return major, minor, patch, nil
}
