package measures

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud/pagination"
)

// Measure is an datapoint thats is composed with a timestamp and a value.
type Measure struct {
	// TimeStamp represents a timestamp of when measure was pushed into the Gnocchi.
	TimeStamp string `json:"-"`

	// Granularity is a level of precision that is kept when aggregating data.
	Granularity float64 `json:"-"`

	// Value represents a value of data that was pushed into the Gnocchi.
	Value float64 `json:"-"`
}

/*
UnmarshalJSON helps to unmarshal response from reading Gnocchi measures.

Gnocchi APIv1 returns measures in a such format:

[
    [
        "2017-01-08T10:00:00+00:00",
        300.0,
        146.0
    ],
    [
        "2017-01-08T10:05:00+00:00",
        300.0,
        58.0
    ]
]

Helper unmarshals every nested array into the Measure type.
*/
func (r *Measure) UnmarshalJSON(b []byte) error {
	type tmp Measure
	var measuresSlice []interface{}

	var s struct {
		tmp
	}
	err := json.Unmarshal(b, &measuresSlice)
	if err != nil {
		return err
	}

	*r = Measure(s.tmp)
	r.TimeStamp = measuresSlice[0].(string)
	r.Granularity = measuresSlice[1].(float64)
	r.Value = measuresSlice[2].(float64)

	return nil
}

// MeasurePage is the page returned by a pager when traversing over a collection
// of measures.
type MeasurePage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a MeasurePage struct is empty.
func (r MeasurePage) IsEmpty() (bool, error) {
	is, err := ExtractMeasures(r)
	return len(is) == 0, err
}

// ExtractMeasures interprets the results of a single page from a List() call,
// producing a slice of Measures structs.
func ExtractMeasures(r pagination.Page) ([]Measure, error) {
	var s []Measure

	err := (r.(MeasurePage)).ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s, err
}
