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

// MeasureToPush represents a single measure that can be pushed into the Gnocchi API.
type MeasureToPush struct {
	// TimeStamp represents a timestamp of when measure is pushed into the Gnocchi.
	TimeStamp time.Time

	// Value represents a value of data that is pushed into the Gnocchi.
	Value float64
}

// PushOptsBuilder is needed to add measures to the Push request.
type PushOptsBuilder interface {
	ToMeasurePushMap() (map[string]interface{}, error)
}

// PushOpts specifies a parameters for pushing measures into a single metric.
type PushOpts struct {
	// Measures is a set of measures that needs to be pushed to a single metric.
	Measures []MeasureToPush
}

// ToMeasurePushMap constructs a request body from PushOpts.
func (opts PushOpts) ToMeasurePushMap() (map[string]interface{}, error) {
	// Struct measureToPush represents internal MeasureToPush variant with string timestamps.
	type measureToPush struct {
		TimeStamp string  `json:"timestamp"`
		Value     float64 `json:"value"`
	}
	type pushOpts struct {
		Measures []measureToPush
	}

	// Convert exported PushOpts to internal pushOpts variant that contains measures with string timestamps.
	internalMeasures := make([]measureToPush, len(opts.Measures))
	for i, m := range opts.Measures {
		internalMeasures[i] = measureToPush{
			TimeStamp: m.TimeStamp.Format(gnocchi.RFC3339NanoNoTimezone),
			Value:     m.Value,
		}
	}
	internalPushOpts := pushOpts{
		Measures: internalMeasures,
	}

	b, err := gophercloud.BuildRequestBody(internalPushOpts, "")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Push requests the creation of a new measures in the single Gnocchi metric.
func Push(client *gophercloud.ServiceClient, metricID string, opts PushOptsBuilder) (r PushResult) {
	b, err := opts.ToMeasurePushMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(pushURL(client, metricID), b["Measures"], &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
		MoreHeaders: map[string]string{
			"Accept": "application/json, */*",
		},
	})
	return
}
