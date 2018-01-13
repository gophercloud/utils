package testing

import (
	"github.com/gophercloud/utils/gnocchi/metric/v1/archivepolicies"
	"github.com/gophercloud/utils/gnocchi/metric/v1/metrics"
)

// MetricsListResult represents a raw server response from a server to a list call.
const MetricsListResult = `[
    {
        "archive_policy": {
        "aggregation_methods": [
            "max",
            "min"
        ],
        "back_window": 0,
        "definition": [
            {
                "granularity": "1:00:00",
                "points": 2304,
                "timespan": "96 days, 0:00:00"
            },
            {
                "granularity": "0:05:00",
                "points": 9216,
                "timespan": "32 days, 0:00:00"
            },
            {
                "granularity": "1 day, 0:00:00",
                "points": 400,
                "timespan": "400 days, 0:00:00"
            }
        ],
        "name": "precise"
        },
        "created_by_project_id": "e9dc821ca664406e981820a477e9a761",
        "created_by_user_id": "a23c5b98d42d4df3b961e54d5167eb6d",
        "creator": "a23c5b98d42d4df3b961e54d5167eb6d:e9dc821ca664406e981820a477e9a761",
        "id": "777a01d6-4694-49cb-b86a-5ba9fd4e609e",
        "name": "memory.usage",
        "resource_id": "1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc",
        "unit": "MB"
    },
    {
        "archive_policy": {
            "aggregation_methods": [
                "mean",
                "sum"
            ],
            "back_window": 12,
            "definition": [
                {
                    "granularity": "1:00:00",
                    "points": 2160,
                    "timespan": "90 days, 0:00:00"
                },
                {
                    "granularity": "1 day, 0:00:00",
                    "points": 200,
                    "timespan": "200 days, 0:00:00"
                }
            ],
            "name": "not_so_precise"
        },
        "created_by_project_id": "c6b68a6b413648b0a0eb191bf3222f4d",
        "created_by_user_id": "cb072aacdb494419aeeba5f1c62d1a65",
        "creator": "cb072aacdb494419aeeba5f1c62d1a65:c6b68a6b413648b0a0eb191bf3222f4d",
        "id": "6dbc97c5-bfdf-47a2-b184-02e7fa348d21",
        "name": "cpu.delta",
        "resource_id": "c5dc0c47-f43c-425c-a82f-44d61ee91175",
        "unit": "ns"
    }
]`

// Metric1 is an expected representation of a first metric from the MetricsListResult.
var Metric1 = metrics.Metric{
	ArchivePolicy: archivepolicies.ArchivePolicy{
		AggregationMethods: []string{
			"max",
			"min",
		},
		BackWindow: 0,
		Definition: []archivepolicies.ArchivePolicyDefinition{
			{
				Granularity: "1:00:00",
				Points:      2304,
				TimeSpan:    "96 days, 0:00:00",
			},
			{
				Granularity: "0:05:00",
				Points:      9216,
				TimeSpan:    "32 days, 0:00:00",
			},
			{
				Granularity: "1 day, 0:00:00",
				Points:      400,
				TimeSpan:    "400 days, 0:00:00",
			},
		},
		Name: "precise",
	},
	CreatedByProjectID: "e9dc821ca664406e981820a477e9a761",
	CreatedByUserID:    "a23c5b98d42d4df3b961e54d5167eb6d",
	Creator:            "a23c5b98d42d4df3b961e54d5167eb6d:e9dc821ca664406e981820a477e9a761",
	ID:                 "777a01d6-4694-49cb-b86a-5ba9fd4e609e",
	Name:               "memory.usage",
	ResourceID:         "1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc",
	Unit:               "MB",
}

// Metric2 is an expected representation of a second metric from the MetricsListResult.
var Metric2 = metrics.Metric{
	ArchivePolicy: archivepolicies.ArchivePolicy{
		AggregationMethods: []string{
			"mean",
			"sum",
		},
		BackWindow: 12,
		Definition: []archivepolicies.ArchivePolicyDefinition{
			{
				Granularity: "1:00:00",
				Points:      2160,
				TimeSpan:    "90 days, 0:00:00",
			},
			{
				Granularity: "1 day, 0:00:00",
				Points:      200,
				TimeSpan:    "200 days, 0:00:00",
			},
		},
		Name: "not_so_precise",
	},
	CreatedByProjectID: "c6b68a6b413648b0a0eb191bf3222f4d",
	CreatedByUserID:    "cb072aacdb494419aeeba5f1c62d1a65",
	Creator:            "cb072aacdb494419aeeba5f1c62d1a65:c6b68a6b413648b0a0eb191bf3222f4d",
	ID:                 "6dbc97c5-bfdf-47a2-b184-02e7fa348d21",
	Name:               "cpu.delta",
	ResourceID:         "c5dc0c47-f43c-425c-a82f-44d61ee91175",
	Unit:               "ns",
}
