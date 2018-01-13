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
*/
package metrics
