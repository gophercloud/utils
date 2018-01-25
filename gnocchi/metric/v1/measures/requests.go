package measures

import (
	"net/url"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/gnocchi"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToMeasureListQuery() (string, error)
}

// ListOpts allows to provide additional options to the Gnocchi measures List request.
type ListOpts struct {
	// Refresh can be used to force any unprocessed measures to be handled in the Gnocchi
	// to ensure that List request returns all aggregates.
	Refresh bool `q:"refresh"`

	// Start is a start of time time range for the measures.
	Start *time.Time

	// Stop is a stop of time time range for the measures.
	Stop *time.Time

	// Aggregation is a needed aggregation method for returned measures.
	// Gnocchi returns "mean" by default.
	Aggregation string `q:"aggregation"`

	// Granularity is a needed time between two series of measures to retreive.
	// Gnocchi will response with all granularities for available measures by default.
	Granularity string `q:"granularity"`

	// Resample allows to select different granularity instead of those that were defined in the
	// archive policy.
	Resample string `q:"resample"`
}

// ToMeasureListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToMeasureListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	params := q.Query()

	if opts.Start != nil {
		params.Add("start", opts.Start.Format(gnocchi.RFC3339NanoNoTimezone))
	}

	if opts.Stop != nil {
		params.Add("stop", opts.Stop.Format(gnocchi.RFC3339NanoNoTimezone))
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// measures.
// It accepts a ListOpts struct, which allows you to provide options to a Gnocchi measures List request.
func List(c *gophercloud.ServiceClient, metricID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c, metricID)
	if opts != nil {
		query, err := opts.ToMeasureListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return MeasurePage{pagination.SinglePageBase(r)}
	})
}
