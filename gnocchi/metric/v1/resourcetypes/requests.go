package resourcetypes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List makes a request against the Gnocchi API to list resource types.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return ResourceTypePage{pagination.SinglePageBase(r)}
	})
}

// Get retrieves a specific Gnocchi resource type based on its name.
func Get(c *gophercloud.ServiceClient, resourceTypeName string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, resourceTypeName), &r.Body, nil)
	return
}

// CreateOptsBuilder allows to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToResourceTypeCreateMap() (map[string]interface{}, error)
}

// AttributeOpts represents options of a single resource type attribute that
// can be created in the Gnocchi.
type AttributeOpts struct {
	// Type is an attribute type.
	Type string `json:"type"`

	// Details represents different attribute fields.
	Details map[string]interface{} `json:"-"`
}

// ToMap is a helper function to convert individual AttributeOpts structure into a sub-map.
func (opts AttributeOpts) ToMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	if opts.Details != nil {
		for k, v := range opts.Details {
			b[k] = v
		}
	}
	return b, nil
}

// CreateOpts specifies parameters of a new Gnocchi resource type.
type CreateOpts struct {
	// Attributes is a collection of keys and values of different resource types.
	Attributes map[string]AttributeOpts `json:"-"`

	// Name is a human-readable resource type identifier.
	Name string `json:"name" required:"true"`
}

// ToResourceTypeCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToResourceTypeCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	// Create resource type without attributes if they're omitted.
	if opts.Attributes == nil {
		return b, nil
	}

	attributes := make(map[string]interface{}, len(opts.Attributes))
	for k, v := range opts.Attributes {
		attributesMap, err := v.ToMap()
		if err != nil {
			return nil, err
		}
		attributes[k] = attributesMap
	}

	b["attributes"] = attributes
	return b, nil
}

// Create requests the creation of a new Gnocchi resource type on the server.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToResourceTypeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}
