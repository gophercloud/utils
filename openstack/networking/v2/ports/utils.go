package ports

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

// IDFromName is a convenience function that returns a port's ID given its
// name. Errors when the number of items found is not one.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	IDs, err := IDsFromName(client, name)
	if err != nil {
		return "", err
	}

	switch count := len(IDs); count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "port"}
	case 1:
		return IDs[0], nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "port"}
	}
}

// IDsFromName returns zero or more IDs corresponding to a name. The returned
// error is only non-nil in case of failure.
func IDsFromName(client *gophercloud.ServiceClient, name string) ([]string, error) {
	pages, err := ports.List(client, ports.ListOpts{
		Name: name,
	}).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ports.ExtractPorts(pages)
	if err != nil {
		return nil, err
	}

	IDs := make([]string, len(all))
	for i := range all {
		IDs[i] = all[i].ID
	}

	return IDs, nil
}
