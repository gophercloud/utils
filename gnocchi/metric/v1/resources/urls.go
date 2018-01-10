package resources

import "github.com/gophercloud/gophercloud"

const resourcePath = "resource"

func rootURL(c *gophercloud.ServiceClient, resourceType string) string {
	if resourceType == "" {
		resourceType = "generic"
	}
	return c.ServiceURL(resourcePath, resourceType)
}

func resourceURL(c *gophercloud.ServiceClient, resourceType, resourceID string) string {
	if resourceType == "" {
		resourceType = "generic"
	}
	return c.ServiceURL(resourcePath, resourceType, resourceID)
}

func listURL(c *gophercloud.ServiceClient, resourceType string) string {
	return rootURL(c, resourceType)
}

func getURL(c *gophercloud.ServiceClient, resourceType, resourceID string) string {
	return resourceURL(c, resourceType, resourceID)
}
