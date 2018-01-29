package v1

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
	"github.com/satori/go.uuid"
)

// CreateGenericResource will create a Gnocchi resource with a generic type.
// An error will be returned if the resource could not be created.
func CreateGenericResource(t *testing.T, client *gophercloud.ServiceClient) (*resources.Resource, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	randomDay := tools.RandomInt(1, 100)
	now := time.Now().UTC().AddDate(0, 0, -randomDay)
	createOpts := resources.CreateOpts{
		ID:        id.String(),
		StartedAt: &now,
		Metrics: map[string]interface{}{
			"cpu.delta": map[string]string{
				"archive_policy_name": "medium",
			},
		},
	}
	resourceType := "generic"
	t.Logf("Attempting to create a generic Gnocchi resource")

	resource, err := resources.Create(client, resourceType, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the generic Gnocchi resource.")
	return resource, nil
}

// DeleteResource will delete a Gnocchi resource with specified type and ID.
// A fatal error will occur if the delete was not successful.
func DeleteResource(t *testing.T, client *gophercloud.ServiceClient, resourceType, resourceID string) {
	t.Logf("Attempting to delete the Gnocchi resource: %s", resourceID)

	err := resources.Delete(client, resourceType, resourceID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete the Gnocchi resource %s: %v", resourceID, err)
	}

	t.Logf("Deleted the Gnocchi resource: %s", resourceID)
}
