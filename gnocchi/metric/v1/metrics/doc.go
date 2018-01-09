/*
Package metrics provides the ability to retrieve metrics through the Gnocchi API.

Example of Listing metrics

	listOpts := metrics.ListOpts{
		Limit: 25,
	}

	allPages, err := metrics.List(gnocchiClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allMetrics, err := metrics.ExtractMetrics(allPages)
	if err != nil {
		panic(err)
	}

	for _, metric := range allMetrics {
		fmt.Printf("%+v\n", metric)
	}

Example of Getting a metric

	metricID = "9e5a6441-1044-4181-b66e-34e180753040"
	metric, err := metrics.Get(gnocchiClient, metricID).Extract()
	if err != nil {
		panic(err)
	}
*/
package metrics
