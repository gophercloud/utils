package measures

import "github.com/gophercloud/gophercloud"

const (
	resourcePath     = "metric"
	batchMetricsPath = "batch/metrics"
)

func resourceURL(c *gophercloud.ServiceClient, metricID string) string {
	return c.ServiceURL(resourcePath, metricID, "measures")
}

func listURL(c *gophercloud.ServiceClient, metricID string) string {
	return resourceURL(c, metricID)
}

func createURL(c *gophercloud.ServiceClient, metricID string) string {
	return resourceURL(c, metricID)
}

func createBatchMetricsURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(batchMetricsPath, "measures")
}
