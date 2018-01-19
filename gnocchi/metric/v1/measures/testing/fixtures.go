package testing

import "github.com/gophercloud/utils/gnocchi/metric/v1/measures"

// MeasuresListResult represents a raw server response from a server to a List call.
const MeasuresListResult = `
[
    [
        "2018-01-10T12:00:00+00:00",
        3600.0,
        15.0
    ],
    [
        "2018-01-10T13:00:00+00:00",
        3600.0,
        10.0
    ],
    [
        "2018-01-10T14:00:00+00:00",
        3600.0,
        20.0
    ]
]
`

// ListArchivePoliciesExpected represents an expected repsonse from a List request.
var ListArchivePoliciesExpected = []measures.Measure{
	{
		TimeStamp:   "2018-01-10T12:00:00+00:00",
		Granularity: 3600.0,
		Value:       15.0,
	},
	{
		TimeStamp:   "2018-01-10T13:00:00+00:00",
		Granularity: 3600.0,
		Value:       10.0,
	},
	{
		TimeStamp:   "2018-01-10T14:00:00+00:00",
		Granularity: 3600.0,
		Value:       20.0,
	},
}
