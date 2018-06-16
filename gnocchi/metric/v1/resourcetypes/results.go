package resourcetypes

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a Gnocchi resource type.
func (r commonResult) Extract() (*ResourceType, error) {
	var s *ResourceType
	err := r.ExtractInto(&s)
	return s, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Gnocchi resource type.
type GetResult struct {
	commonResult
}

// ResourceType represents custom Gnocchi resource type.
type ResourceType struct {
	// Attributes is a collection of keys and values of different resource types.
	Attributes []ResourceTypeAttribute `json:"-"`

	// Name is a human-readable resource type identifier.
	Name string `json:"name"`

	// State represents current status of a resource type.
	State string `json:"state"`
}

// ResourceTypeAttribute represents single attribute of a Gnocchi resource type.
type ResourceTypeAttribute struct {
	// Name a human-readable attribute identifier.
	Name string `json:"-"`

	// MaxLength contains maximum length of an attribute value.
	MaxLength int `json:"max_length"`

	// MinLength contains minimum length of an attribute value.
	MinLength int `json:"min_length"`

	// Required shows if that attribute is required.
	Required bool `json:"required"`

	// Type is an attribute type.
	Type string `json:"type"`
}

// UnmarshalJSON helps to unmarshal ResourceType fields into needed values.
func (r *ResourceType) UnmarshalJSON(b []byte) error {
	type tmp ResourceType
	var s struct {
		tmp
		Attributes map[string]interface{} `json:"attributes"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ResourceType(s.tmp)

	if s.Attributes == nil {
		return nil
	}

	// Populate attributes from the JSON map structure.
	attributes := make([]ResourceTypeAttribute, len(s.Attributes))
	idx := 0
	for attributeName, attributeValues := range s.Attributes {
		attributes[idx] = ResourceTypeAttribute{
			Name: attributeName,
		}

		var attributeValuesMap map[string]interface{}
		var ok bool

		if attributeValuesMap, ok = attributeValues.(map[string]interface{}); !ok {
			// Got some strange resource type attribute representation, skip it.
			continue
		}

		// Populate attribute values, check types and skip invalid ones.
		for k, v := range attributeValuesMap {
			switch {
			case k == "max_length":
				if maxLength, ok := v.(float64); ok {
					attributes[idx].MaxLength = int(maxLength)
				}
			case k == "min_length":
				if minLength, ok := v.(float64); ok {
					attributes[idx].MinLength = int(minLength)
				}
			case k == "required":
				if required, ok := v.(bool); ok {
					attributes[idx].Required = required
				}
			case k == "type":
				if attributeType, ok := v.(string); ok {
					attributes[idx].Type = attributeType
				}
			}
		}
		idx++
	}

	r.Attributes = attributes

	return err
}

// ResourceTypePage abstracts the raw results of making a List() request against
// the Gnocchi API.
//
// As Gnocchi API may freely alter the response bodies of structures
// returned to the client, you may only safely access the data provided through
// the ExtractResources call.
type ResourceTypePage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a ResourceTypePage struct is empty.
func (r ResourceTypePage) IsEmpty() (bool, error) {
	is, err := ExtractResourceTypes(r)
	return len(is) == 0, err
}

// ExtractResourceTypes interprets the results of a single page from a List() call,
// producing a slice of ResourceType structs.
func ExtractResourceTypes(r pagination.Page) ([]ResourceType, error) {
	var s []ResourceType
	err := (r.(ResourceTypePage)).ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s, err
}
