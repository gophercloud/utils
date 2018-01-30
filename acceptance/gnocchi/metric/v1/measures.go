package v1

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/gnocchi/metric/v1/measures"
)

// CreateMeasures will create measures inside a single Gnocchi metric. An error will be returned if the
// measures could not be created.
func CreateMeasures(t *testing.T, client *gophercloud.ServiceClient, metricID string) error {
	currentTimestamp := time.Now().UTC()
	pastHourTimestamp := currentTimestamp.Add(-1 * time.Hour)
	currentValue := float64(tools.RandomInt(100, 200))
	pastHourValue := float64(tools.RandomInt(500, 600))
	measuresToCreate := []measures.MeasureOpts{
		{
			Timestamp: &currentTimestamp,
			Value:     currentValue,
		},
		{
			Timestamp: &pastHourTimestamp,
			Value:     pastHourValue,
		},
	}
	createOpts := measures.CreateOpts{
		Measures: measuresToCreate,
	}

	t.Logf("Attempting to create measures inside a Gnocchi metric %s", metricID)

	if err := measures.Create(client, metricID, createOpts).ExtractErr(); err != nil && err.Error() != "EOF" {
		return err
	}

	t.Logf("Successfully created measures inside the Gnocchi metric %s", metricID)
	return nil
}

// BatchMetricsMeasures will create measures inside different metrics via batch request. An error will be returned if the
// measures could not be created.
func BatchMetricsMeasures(t *testing.T, client *gophercloud.ServiceClient, metricIDs ...string) error {
	currentTimestamp := time.Now().UTC()
	pastHourTimestamp := currentTimestamp.Add(-1 * time.Hour)
	currentValue := float64(tools.RandomInt(100, 200))
	pastHourValue := float64(tools.RandomInt(500, 600))
	batchOpts := make(map[string][]measures.MeasureOpts)
	for _, m := range metricIDs {
		batchOpts[m] = []measures.MeasureOpts{
			{
				Timestamp: &currentTimestamp,
				Value:     currentValue,
			},
			{
				Timestamp: &pastHourTimestamp,
				Value:     pastHourValue,
			},
		}
	}
	createOpts := measures.BatchMetricsOpts{
		BatchOpts: batchOpts,
	}

	t.Logf("Attempting to create measures inside Gnocchi metrics via batch request")

	if err := measures.BatchMetrics(client, createOpts).ExtractErr(); err != nil && err.Error() != "EOF" {
		return err
	}

	t.Logf("Successfully created measures inside Gnocchi metrics")
	return nil
}
