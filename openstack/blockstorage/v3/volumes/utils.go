package volumes

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

// IDFromName is a convenience function that returns a volume's ID given its
// name. Errors when the number of items found is not one.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	IDs, err := IDsFromName(client, name)
	if err != nil {
		return "", err
	}

	switch count := len(IDs); count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "volume"}
	case 1:
		return IDs[0], nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "volume"}
	}
}

// IDsFromName returns zero or more IDs corresponding to a name. The returned
// error is only non-nil in case of failure.
func IDsFromName(client *gophercloud.ServiceClient, name string) ([]string, error) {
	pages, err := volumes.List(client, volumes.ListOpts{
		Name: name,
	}).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := volumes.ExtractVolumes(pages)
	if err != nil {
		return nil, err
	}

	IDs := make([]string, len(all))
	for i := range all {
		IDs[i] = all[i].ID
	}

	return IDs, nil
}

// SafeVolumeDelete is a convenience function that returns an error if a volume can't be removed because it's in a wrong
// state (e.g. it's attached to a server).
func SafeVolumeDelete(client *gophercloud.ServiceClient, id string, opts volumes.DeleteOpts) error {
	vol, err := volumes.Get(client, id).Extract()
	if err != nil {
		return err
	}
	status := vol.Status
	switch status {
	// https://docs.openstack.org/api-ref/block-storage/v3/index.html?expanded=delete-a-volume-detail#delete-a-volume
	case "available", "error", "error_restoring", "error_extending", "error_managing":
		return volumes.Delete(client, id, opts).ExtractErr()
	default:
		return fmt.Errorf("volume %s is in %s state, can't delete", id, status)
	}
}
