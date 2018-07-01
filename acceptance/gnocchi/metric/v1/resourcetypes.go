package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/utils/gnocchi/metric/v1/resourcetypes"
)

// CreateResourceType creates Gnocchi resource type. An error will be returned if the
// resource type could not be created.
func CreateResourceType(t *testing.T, client *gophercloud.ServiceClient) (*resourcetypes.ResourceType, error) {
	resourceTypeName := tools.RandomString("TESTACCT-", 8)
	attributeStringName := tools.RandomString("TESTACCT-ATTRIBUTE-", 8)
	attributeUUIDName := tools.RandomString("TESTACCT-ATTRIBUTE-", 8)

	createOpts := resourcetypes.CreateOpts{
		Name: resourceTypeName,
		Attributes: map[string]resourcetypes.AttributeOpts{
			attributeStringName: resourcetypes.AttributeOpts{
				Type: "string",
				Details: map[string]interface{}{
					"max_length": 128,
					"required":   false,
				},
			},
			attributeUUIDName: resourcetypes.AttributeOpts{
				Type: "uuid",
				Details: map[string]interface{}{
					"required": true,
				},
			},
		},
	}
	t.Logf("Attempting to create a Gnocchi resource type")

	resourceType, err := resourcetypes.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the Gnocchi resource type.")
	return resourceType, nil
}
