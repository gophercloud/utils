package testing

import "github.com/gophercloud/utils/gnocchi/metric/v1/resourcetypes"

// ResourceTypeListResult represents raw server response from a server to a list call.
const ResourceTypeListResult = `[
    {
        "attributes": {},
        "name": "generic",
        "state": "active"
    },
    {
        "attributes": {
            "parent_id": {
                "required": false,
                "type": "uuid"
            }
        },
        "name": "identity_project",
        "state": "active"
    },
    {
        "attributes": {
            "host": {
                "max_length": 128,
                "min_length": 0,
                "required": true,
                "type": "string"
            }
        },
        "name": "compute_instance",
        "state": "active"
    }
]`

// ResourceType1 is an expected representation of a first resource from the ResourceTypeListResult.
var ResourceType1 = resourcetypes.ResourceType{
	Name:       "generic",
	State:      "active",
	Attributes: map[string]resourcetypes.Attribute{},
}

// ResourceType2 is an expected representation of a first resource from the ResourceTypeListResult.
var ResourceType2 = resourcetypes.ResourceType{
	Name:  "identity_project",
	State: "active",
	Attributes: map[string]resourcetypes.Attribute{
		"parent_id": resourcetypes.Attribute{
			Type: "uuid",
			Details: map[string]interface{}{
				"required": false,
			},
		},
	},
}

// ResourceType3 is an expected representation of a first resource from the ResourceTypeListResult.
var ResourceType3 = resourcetypes.ResourceType{
	Name:  "compute_instance",
	State: "active",
	Attributes: map[string]resourcetypes.Attribute{
		"host": resourcetypes.Attribute{
			Type: "string",
			Details: map[string]interface{}{
				"max_length": float64(128),
				"min_length": float64(0),
				"required":   true,
			},
		},
	},
}
