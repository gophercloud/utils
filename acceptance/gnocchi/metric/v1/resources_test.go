// +build acceptance metric resources

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/acceptance/clients"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
)

func TestResourcesCRUD(t *testing.T) {
	client, err := clients.NewGnocchiV1Client()
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi client: %v", err)
	}

	genericResource, err := CreateGenericResource(t, client)
	if err != nil {
		t.Fatalf("Unable to create a generic Gnocchi resource: %v", err)
	}

	tools.PrintResource(t, genericResource)
}

func TestResourcesList(t *testing.T) {
	client, err := clients.NewGnocchiV1Client()
	if err != nil {
		t.Fatalf("Unable to create a Gnocchi client: %v", err)
	}

	opts := resources.ListOpts{}
	resourceType := "generic"
	allPages, err := resources.List(client, opts, resourceType).AllPages()
	if err != nil {
		t.Fatalf("Unable to list resources: %v", err)
	}

	allResources, err := resources.ExtractResources(allPages)
	if err != nil {
		t.Fatalf("Unable to extract resources: %v", err)
	}

	for _, resource := range allResources {
		tools.PrintResource(t, resource)
	}
}
