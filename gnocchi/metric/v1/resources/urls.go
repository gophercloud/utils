package resources

import "github.com/gophercloud/gophercloud"

const resourcePath = "resource"

func rootURL(c *gophercloud.ServiceClient, resourceType string) string {
	if resourceType == "" {
		resourceType = "generic"
	}
	return c.ServiceURL(resourcePath, resourceType)
}

func listURL(c *gophercloud.ServiceClient, resourceType string) string {
	return rootURL(c, resourceType)
}

func getURL(c *gophercloud.ServiceClient, resourceID string, resourceType string) string {
	if resourceType == "" {
		resourceType = "generic"
	}
	return c.ServiceURL(resourcePath, resourceType, resourceID)
}
