package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// IDFromName is a convienience function that returns an image's ID given its
// name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	allPages, err := images.List(client, images.ListOpts{
		Name: name,
	}).AllPages()
	if err != nil {
		return "", err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return "", err
	}

	switch len(allImages) {
	case 0:
		err := &gophercloud.ErrResourceNotFound{}
		err.ResourceType = "image"
		err.Name = name
		return "", err
	case 1:
		return allImages[0].ID, nil
	default:
		err := &gophercloud.ErrMultipleResourcesFound{}
		err.ResourceType = "image"
		err.Name = name
		err.Count = len(allImages)
		return "", err
	}
}
