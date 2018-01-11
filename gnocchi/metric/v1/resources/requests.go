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

// CreateOptsBuilder allows to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToResourceCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies parameters of a new Gnocchi resource.
type CreateOpts struct {
	// CreatedByProjectID contains the id of the Identity project that
	// was used for a resource creation.
	CreatedByProjectID string `json:"created_by_project_id,omitempty"`

	// CreatedByUserID contains the id of the Identity user
	// that created the Gnocchi resource.
	CreatedByUserID string `json:"created_by_user_id,omitempty"`

	// Creator shows who created the resource.
	// Usually it contains concatenated string with values from
	// "created_by_user_id" and "created_by_project_id" fields.
	Creator string `json:"creator,omitempty"`

	// Metrics field can be used to link existing metrics in the resource
	// or to create metrics with the resource at the same time to save
	// some requests.
	Metrics map[string]interface{} `json:"metrics,omitempty"`

	// ID uniquely identifies the Gnocchi resource.
	ID string `json:"id,omitempty"`

	// OriginalResourceID is the orginal resource id. It can be different from the
	// regular ID field.
	OriginalResourceID string `json:"original_resource_id,omitempty"`

	// ProjectID is the Identity project of the resource.
	ProjectID string `json:"project_id,omitempty"`

	// UserID is the Identity user of the resource.
	UserID string `json:"user_id,omitempty"`

	// StartedAt is a resource creation timestamp.
	StartedAt string `json:"started_at,omitempty"`

	// EndedAt is a timestamp of when the resource has ended.
	EndedAt string `json:"ended_at,omitempty"`
}

// ToResourceCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToResourceCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create requests the creation of a new Gnocchi resource on the server.
func Create(client *gophercloud.ServiceClient, resourceType string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToResourceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client, resourceType), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	return
}
