package resources

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToResourceListQuery() (string, error)
}

// ListOpts allows the limiting and sorting of paginated collections through
// the Gnocchi API.
type ListOpts struct {
	// Details allows to list resources with all attributes.
	Details bool `q:"details"`

	// Limit allows to limits count of resources in the response.
	Limit int `q:"limit"`

	// Marker is used for pagination.
	Marker string `q:"marker"`

	// SortKey allows to sort resources in the response by key.
	SortKey string `q:"sort_key"`

	// SortDir allows to set the direction of sorting.
	// Can be `asc` or `desc`.
	SortDir string `q:"sort_dir"`
}

// ToResourceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToResourceListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// resources. It accepts a ListOpts struct, which allows you to limit and sort
// the returned collection for a greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder, resourceType string) pagination.Pager {
	url := listURL(c, resourceType)
	if opts != nil {
		query, err := opts.ToResourceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ResourcePage{pagination.SinglePageBase(r)}
	})
}

// Get retrieves a specific Gnocchi resource based on its type and ID.
func Get(c *gophercloud.ServiceClient, resourceType string, resourceID string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, resourceType, resourceID), &r.Body, nil)
	return
}
