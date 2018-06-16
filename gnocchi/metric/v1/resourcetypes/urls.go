package resourcetypes

import "github.com/gophercloud/gophercloud"

const resourcePath = "resource_type"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
