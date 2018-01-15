package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
)

// CreateResource will create Gnocchi resource. An error will be returned if the
// resource could not be created.
func CreateResource(t *testing.T, client *gophercloud.ServiceClient) (*resources.Resource, error) {
	createOpts := resources.CreateOpts{
		ID: "00000000-0000-dead-beef-111111111111",
	}
	resourceType := "generic"
	t.Logf("Attempting to create a Gnocchi resource")

	resource, err := resources.Create(client, resourceType, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the Gnocchi resource.")
	return resource, nil
}
