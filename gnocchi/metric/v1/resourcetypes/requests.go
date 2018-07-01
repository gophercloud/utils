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
