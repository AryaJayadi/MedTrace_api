package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

// OptionalTime is a wrapper around time.Time for custom JSON marshalling.
// It marshals to an empty string if the time.Time value is zero.
type OptionalTime struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface for OptionalTime.
// If the time.Time value is its zero value, it marshals as an empty JSON string ("").
// Otherwise, it marshals using the default time.Time format (RFC3339Nano).
func (ot OptionalTime) MarshalJSON() ([]byte, error) {
	if ot.Time.IsZero() {
		return []byte(`""`), nil // Marshal as an empty JSON string
	}
	// Use standard marshalling for non-zero time.
	// Need to explicitly marshal to avoid infinite recursion if json.Marshal(ot.Time) was used directly
	// on a type that also has custom marshalling. Standard time.Time does not, so this is safe.
	return json.Marshal(ot.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface for OptionalTime.
// It handles unmarshalling from a JSON string. If the string is empty or
// represents Go's zero time ("0001-01-01T00:00:00Z"), it sets the time.Time to its zero value.
func (ot *OptionalTime) UnmarshalJSON(data []byte) error {
	// First, try to unmarshal as a string.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s == "" || s == "0001-01-01T00:00:00Z" {
			ot.Time = time.Time{} // Set to zero value
			return nil
		}
		// If not empty or zero string, try to parse it as a standard time.
		// Go's default time.Time unmarshalling handles RFC3339Nano.
		// We can replicate that or be more specific.
		parsedTime, errParse := time.Parse(time.RFC3339Nano, s)
		if errParse != nil {
			// Fallback to RFC3339 if Nano parsing fails
			parsedTime, errParse = time.Parse(time.RFC3339, s)
			if errParse != nil {
				return fmt.Errorf("OptionalTime: unable to parse date string '%s': %w", s, errParse)
			}
		}
		ot.Time = parsedTime
		return nil
	}

	// If it wasn't a string (e.g., `null` or some other JSON type),
	// try to unmarshal directly into a time.Time.
	// If `data` is `null`, this will set ot.Time to its zero value without error.
	var t time.Time
	if errDirect := json.Unmarshal(data, &t); errDirect != nil {
		// This path is less common if the source JSON always uses strings for dates,
		// but it provides a fallback.
		return fmt.Errorf("OptionalTime: data '%s' is not a recognized time string or null: %w", string(data), errDirect)
	}
	ot.Time = t
	return nil
}
