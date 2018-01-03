package metrics

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/gnocchi/metric/v1/archivepolicies"
)

type metricResult struct {
	gophercloud.Result
}

// Metric is an entity storing aggregates identified by an UUID.
// It can be attached to a resource using a name.
// How a metric stores its aggregates is defined by the archive policy
// it is associated to.
type Metric struct {
	// ArchivePolicy is a Gnocchi archive policy that describes the aggregate
	// storage policy of a metric.
	ArchivePolicy archivepolicies.ArchivePolicy `json:"archive_policy"`

	// CreatedByProjectID contains the id of the Identity project that
	// was used for a metric creation.
	CreatedByProjectID string `json:"created_by_project_id"`

	// CreatedByUserID contains the id of the Identity user
	// that created the Gnocchi metric.
	CreatedByUserID string `json:"created_by_user_id"`

	// Creator shows who created the metric.
	// Usually it contains concatenated string with values from
	// "created_by_user_id" and "created_by_project_id" fields.
	Creator string `json:"creator"`

	// ID uniquely identifies the Gnocchi metric.
	ID string `json:"id"`

	// Name is a human-readable name for the Gnocchi metric.
	Name string `json:"name"`

	// ResourceID identifies the associated Gnocchi resource of the metric.
	ResourceID string `json:"resource_id"`

	// Unit is a unit of measurement for measures of that Gnocchi metric.
	Unit string `json:"unit"`
}

// MetricPage is the page returned by a pager when traversing over a collection
// of metrics.
type MetricPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a MetricPage struct is empty.
func (r MetricPage) IsEmpty() (bool, error) {
	is, err := ExtractMetrics(r)
	return len(is) == 0, err
}

// ExtractMetrics interprets the results of a single page from a List() call,
// producing a slice of Metric structs.
func ExtractMetrics(r pagination.Page) ([]Metric, error) {
	var s []Metric
	err := (r.(MetricPage)).ExtractInto(&s)
	if err != nil {
		return nil, err
	}

	return s, err
}
