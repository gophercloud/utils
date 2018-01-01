package resources

import "github.com/gophercloud/gophercloud"

const resourcePath = "resource"

func listURL(c *gophercloud.ServiceClient, resourceType string) string {
	if resourceType == "" {
		resourceType = "generic"
	}
	return c.ServiceURL(resourcePath, resourceType)
}
