package enumUtils

import (
	"fmt"
	"strings"
)

// Enum represents an enumeration type with predefined options.
type Enum struct {
	Allowed []string // Allowed contains the list of allowed options.
	Value   string   // Value represents the current value of the enumeration.
}

// NewEnum creates a new Enum object with the given list of allowed options and default value.
func NewEnum(allowed []string, d string) *Enum {
	return &Enum{
		Allowed: allowed,
		Value:   d,
	}
}

// String returns the string representation of the current value.
func (a Enum) String() string {
	return a.Value
}

// Set sets the value of the enumeration.
// It returns an error if the provided value is not in the list of allowed options.
func (a *Enum) Set(p string) error {
	// isIncluded checks if a value exists in a slice.
	isIncluded := func(opts []string, val string) bool {
		for _, opt := range opts {
			if val == opt {
				return true
			}
		}
		return false
	}

	// Check if the provided value is in the list of allowed options.
	if !isIncluded(a.Allowed, p) {
		return fmt.Errorf("%s is not included in %s", p, strings.Join(a.Allowed, ","))
	}

	// Set the value of the enumeration.
	a.Value = p
	return nil
}

// Type returns the type of the enumeration.
func (a *Enum) Type() string {
	return "string"
}
