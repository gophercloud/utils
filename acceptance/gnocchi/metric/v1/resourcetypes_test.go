// +build acceptance metric resourcetypes

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/acceptance/clients"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resourcetypes"
)

func TestResourceTypesList(t *testing.T) {
	client, err := clients.NewGnocchiV1Client()
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi client: %v", err)
	}

	allPages, err := resourcetypes.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to list resource types: %v", err)
	}

	allResourceTypes, err := resourcetypes.ExtractResourceTypes(allPages)
	if err != nil {
		t.Fatalf("Unable to extract resource types: %v", err)
	}

	for _, resourceType := range allResourceTypes {
		tools.PrintResource(t, resourceType)
	}
}
