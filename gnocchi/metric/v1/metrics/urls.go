package metrics

import "github.com/gophercloud/gophercloud"

const resourcePath = "metric"

func resourceURL(c *gophercloud.ServiceClient, metricID string) string {
	return c.ServiceURL(resourcePath, metricID)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gophercloud.ServiceClient, metricID string) string {
	return resourceURL(c, metricID)
}
