package gnocchi

import (
	"bytes"
	"encoding/json"
	"time"
)

// RFC3339NanoTimezone describes a common timestamp format used by Gnocchi API responses.
const RFC3339NanoTimezone = "2006-01-02T15:04:05.999999+00:00"

// JSONRFC3339NanoTimezone is a type for Gnocchi timestamps.
type JSONRFC3339NanoTimezone time.Time

// UnmarshalJSON helps to unmarshal timestamps from Gnocchi responses to the
// JSONRFC3339NanoTimezone type.
func (jt *JSONRFC3339NanoTimezone) UnmarshalJSON(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	t, err := time.Parse(RFC3339NanoTimezone, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339NanoTimezone(t)
	return nil
}
