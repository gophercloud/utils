package metrics

import "github.com/gophercloud/gophercloud"

const resourcePath = "metric"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}
