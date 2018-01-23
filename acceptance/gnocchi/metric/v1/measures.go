package v1

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/gnocchi/metric/v1/measures"
)

// CreateMeasures will create measures inside a single Gnocchi metric. An error will be returned if the
// measures could not be created.
func CreateMeasures(t *testing.T, client *gophercloud.ServiceClient, metricID string) error {
	currentTimeStamp := time.Now().UTC()
	pastHourTimeStamp := currentTimeStamp.Add(-1 * time.Hour)
	measuresToCreate := []measures.MeasureOpts{
		{
			TimeStamp: currentTimeStamp,
			Value:     100.5,
		},
		{
			TimeStamp: pastHourTimeStamp,
			Value:     500,
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
