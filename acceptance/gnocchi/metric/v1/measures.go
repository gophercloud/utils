package v1

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/gnocchi/metric/v1/measures"
)

// PushMeasures will push measures into a single Gnocchi metric. An error will be returned if the
// measures could not be created.
func PushMeasures(t *testing.T, client *gophercloud.ServiceClient, metricID string) error {
	currentTimeStamp := time.Now().UTC()
	pastHourTimeStamp := currentTimeStamp.Add(-1 * time.Hour)
	measuresToPush := []measures.MeasureToPush{
		{
			TimeStamp: currentTimeStamp,
			Value:     100.5,
		},
		{
			TimeStamp: pastHourTimeStamp,
			Value:     500,
		},
	}
	pushOpts := measures.PushOpts{
		Measures: measuresToPush,
	}

	t.Logf("Attempting to push measures into a Gnocchi metric %s", metricID)

	if err := measures.Push(client, metricID, pushOpts).ExtractErr(); err != nil && err.Error() != "EOF" {
		return err
	}

	t.Logf("Successfully pushed measures into the Gnocchi metric %s", metricID)
	return nil
}
