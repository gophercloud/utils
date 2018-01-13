package resources

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/gnocchi"
)

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
	Metrics map[string]string `json:"metrics"`

	// OriginalResourceID is the orginal resource id. It can be different from the
	// regular ID field.
	OriginalResourceID string `json:"original_resource_id"`

	// ProjectID is the Identity project of the resource.
	ProjectID string `json:"project_id"`

	// RevisionStart is a staring timestamp of the current resource revision.
	RevisionStart time.Time `json:"-"`

	// RevisionEnd is an ending timestamp of the last resource revision.
	RevisionEnd time.Time `json:"-"`

	// StartedAt is a resource creation timestamp.
	StartedAt time.Time `json:"-"`

	// EndedAt is a timestamp of when the resource has ended.
	EndedAt time.Time `json:"-"`

	// Type is a type of the resource.
	Type string `json:"type"`

	// UserID is the Identity user of the resource.
	UserID string `json:"user_id"`
}

// UnmarshalJSON helps to unmarshal Resource fields into needed values.
func (r *Resource) UnmarshalJSON(b []byte) error {
	type tmp Resource
	var s struct {
		tmp
		RevisionStart gnocchi.JSONRFC3339NanoTimezone `json:"revision_start"`
		RevisionEnd   gnocchi.JSONRFC3339NanoTimezone `json:"revision_end"`
		StartedAt     gnocchi.JSONRFC3339NanoTimezone `json:"started_at"`
		EndedAt       gnocchi.JSONRFC3339NanoTimezone `json:"ended_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Resource(s.tmp)

	r.RevisionStart = time.Time(s.RevisionStart)
	r.RevisionEnd = time.Time(s.RevisionEnd)
	r.StartedAt = time.Time(s.StartedAt)
	r.EndedAt = time.Time(s.EndedAt)

	return err
}

// ResourcePage abstracts the raw results of making a List() request against
// the Gnocchi API.
//
// As Gnocchi API may freely alter the response bodies of structures
// returned to the client, you may only safely access the data provided through
// the ExtractResources call.
type ResourcePage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a ResourcePage struct is empty.
func (r ResourcePage) IsEmpty() (bool, error) {
	is, err := ExtractResources(r)
	return len(is) == 0, err
}

// ExtractResources interprets the results of a single page from a List() call,
// producing a slice of Resource structs.
func ExtractResources(r pagination.Page) ([]Resource, error) {
	var s []Resource
	err := (r.(ResourcePage)).ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s, err
}
