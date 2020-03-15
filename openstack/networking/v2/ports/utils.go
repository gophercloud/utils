package ports

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

// IDFromName is a convenience function that returns a port's ID,
// given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	listOpts := ports.ListOpts{
		Name: name,
	}

	pages, err := ports.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ports.ExtractPorts(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "port"}
	case 1:
		return id, nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "port"}
	}
}
