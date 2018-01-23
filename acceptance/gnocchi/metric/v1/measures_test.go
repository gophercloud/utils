// +build acceptance metric measures

package v1

import (
	"testing"

	"github.com/gophercloud/utils/acceptance/clients"
	"github.com/gophercloud/utils/gnocchi/metric/v1/measures"
)

func TestMeasuresCRUD(t *testing.T) {
	client, err := clients.NewGnocchiV1Client()
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi client: %v", err)
	}

	metric, err := CreateMetric(t, client)
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi metric: %v", err)
	}
	defer DeleteMetric(t, client, metric.ID)

	if err := CreateMeasures(t, client, metric.ID); err != nil {
		t.Fatalf("Unable to create measures inside the Gnocchi metric: %v", err)
	}

	listOpts := measures.ListOpts{
		Refresh: true,
	}
	allPages, err := measures.List(client, metric.ID, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to list measures of the metric %s: %v", metric.ID, err)
	}

	metricMeasures, err := measures.ExtractMeasures(allPages)
	if err != nil {
		t.Fatalf("Unable to extract measures: %v", metricMeasures)
	}

	t.Log(metricMeasures)
}
