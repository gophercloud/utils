package resources

import "github.com/gophercloud/gophercloud/pagination"

// Resource is an entity representing anything in your infrastructure
// that you will associate metric(s) with.
// It is identified by a unique ID and can contain attributes.
type Resource struct {
	// CreatedByProjectID contains the id of the Identity project that
	// was used for a resource creation.
	CreatedByProjectID string `json:"created_by_project_id"`

	// CreatedByUserID contains the id of the Identity user
	// that created the Gnocchi resource.
	CreatedByUserID string `json:"created_by_user_id"`

	// Creator shows who created the resource.
	// Usually it contains concatenated string with values from
	// "created_by_user_id" and "created_by_project_id" fields.
	Creator string `json:"creator"`

	// ID uniquely identifies the Gnocchi resource.
	ID string `json:"id"`

	// Metrics are entities that store aggregates.
	// They are identified by an UUIDs.
	Metrics []string `json:"metrics"`

	// OriginalResourceID is the orginal resource id. It can be different from the
	// regular ID field.
	OriginalResourceID string `json:"original_resource_id"`

	// ProjectID is the Identity project of the resource.
	ProjectID string `json:"project_id"`

	// RevisionStart is a staring timestamp of the current resource revision.
	RevisionStart string `json:"revision_start"`

	// RevisionEnd is an ending timestamp of the last resource revision.
	RevisionEnd string `json:"revision_end"`

	// StartedAt is a resource creation timestamp.
	StartedAt string `json:"started_at"`

	// Type is a type of the resource.
	Type string `json:"type"`

	// UserID is the Identity user of the resource.
	UserID string `json:"user_id"`
}

// ResourcePage is the page returned by a pager when traversing over a collection
// of resources.
type ResourcePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of resources has reached
// the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
//
// Gnocchi API doesn't provide special JSON with pagination links. Instead it expects
// a new URL with a resource's id as a marker.
func (r ResourcePage) NextPageURL() (string, error) {
	var s []Resource

	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	lastResource := s[len(s)-1]
	lastResourceID := lastResource.ID
	nextPageURL := r.PageResult.URL.String() + "?marker=" + lastResourceID

	return nextPageURL, nil
}

// IsEmpty checks whether a ResourcePage struct is empty.
func (r ResourcePage) IsEmpty() (bool, error) {
	is, err := ExtractResources(r)
	return len(is) == 0, err
}

// ExtractResources accepts a Page struct, specifically a ResourcePage struct,
// and extracts the elements into a slice of Resource structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractResources(r pagination.Page) ([]Resource, error) {
	var s []Resource
	err := (r.(ResourcePage)).ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s, err
}
