package query

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/gophercloud/gophercloud"
)

func New(listOpts interface{}) *ListOpts {
	availableFields := make(map[string]string)
	{
		t := reflect.TypeOf(listOpts)
		for i := 0; i < t.NumField(); i++ {
			if tag := t.Field(i).Tag.Get("q"); tag != "" {
				availableFields[tag] = t.Field(i).Name
			}
		}
	}

	queryString, err := gophercloud.BuildQueryString(listOpts)

	return &ListOpts{
		allowedFields: availableFields,
		query:         queryString.Query(),
		errs:          joinErrors(err),
	}
}

// ListOpts can be used to list multiple resources.
type ListOpts struct {
	allowedFields map[string]string
	query         url.Values
	errs          error
}

// And adds an arbitrary number of permutations of a single property to filter
// in. When a single ListOpts is called multiple times with the same property
// name, the resulting query contains the resulting intersection (AND). Note
// that how these properties are combined in OpenStack depend on the property.
// For example: passing multiple "id" behaves like an OR. Instead, passing
// multiple "tags" will only return resources that have ALL those tags. This
// helper function only combines the parameters in the most straightforward
// way; please refer to the OpenStack documented behaviour to know how these
// parameters are treated.
//
// ListOpts is currently implemented for three Network resources:
//
// * ports
// * networks
// * subnets
func (o *ListOpts) And(property string, values ...interface{}) *ListOpts {
	if existingValues, ok := o.query[property]; ok {
		// There already are values of the same property: we AND them
		// with the new ones. We only keep the values that exist in
		// both `o.query` AND in `values`.

		// First, to avoid nested loops, we build a hashmap with the
		// new values.
		newValuesSet := make(map[string]struct{})
		for _, newValue := range values {
			newValuesSet[fmt.Sprint(newValue)] = struct{}{}
		}

		// intersectedValues is a slice which will contain the values
		// that we want to keep. They will be at most as many as what
		// we already have; that's what we set the slice capacity to.
		intersectedValues := make([]string, 0, len(existingValues))

		// We add each existing value to intersectedValues if and only
		// if it's also present in the new set.
		for _, existingValue := range existingValues {
			if _, ok := newValuesSet[existingValue]; ok {
				intersectedValues = append(intersectedValues, existingValue)
			}
		}
		o.query[property] = intersectedValues
		return o
	}

	if _, ok := o.allowedFields[property]; !ok {
		o.errs = joinErrors(o.errs, fmt.Errorf("invalid property for the filter: %q", property))
		return o
	}

	for _, v := range values {
		o.query.Add(property, fmt.Sprint(v))
	}

	return o
}

func (o ListOpts) String() string {
	return "?" + o.query.Encode()
}

func (o ListOpts) ToPortListQuery() (string, error) {
	return o.String(), o.errs
}

func (o ListOpts) ToNetworkListQuery() (string, error) {
	return o.String(), o.errs
}

func (o ListOpts) ToSubnetListQuery() (string, error) {
	return o.String(), o.errs
}
