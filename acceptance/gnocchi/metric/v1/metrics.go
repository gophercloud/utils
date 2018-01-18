package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/gnocchi/metric/v1/metrics"
)

// CreateMetric will create Gnocchi metric. An error will be returned if the
// metric could not be created.
func CreateMetric(t *testing.T, client *gophercloud.ServiceClient) (*metrics.Metric, error) {
	// Metric will be created assuming that your Gnocchi's indexer installation was configured with
	// the "gnocchi-manage --noskip-archive-policies-creation" command. So Gnocchi has the default policies:
	// "low", "medium", "high", "bool".
	createOpts := metrics.CreateOpts{
		ArchivePolicyName: "low",
	}
	t.Logf("Attempting to create a Gnocchi metric")

	metric, err := metrics.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the Gnocchi metric.")
	return metric, nil
}
