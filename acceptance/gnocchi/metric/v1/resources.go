package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
	"github.com/satori/go.uuid"
)

// CreateResource will create Gnocchi resource. An error will be returned if the
// resource could not be created.
func CreateResource(t *testing.T, client *gophercloud.ServiceClient) (*resources.Resource, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	createOpts := resources.CreateOpts{
		ID: id.String(),
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
