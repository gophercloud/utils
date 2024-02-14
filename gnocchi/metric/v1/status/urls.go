package status

import "github.com/gophercloud/gophercloud/v2"

const resourcePath = "status"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func getURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
