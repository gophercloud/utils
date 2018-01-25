package measures

import "github.com/gophercloud/gophercloud"

const resourcePath = "metric"

func resourceURL(c *gophercloud.ServiceClient, metricID string) string {
	return c.ServiceURL(resourcePath, metricID, "measures")
}

func listURL(c *gophercloud.ServiceClient, metricID string) string {
	return resourceURL(c, metricID)
}
