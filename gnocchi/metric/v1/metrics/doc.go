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

Example of Creating a metric and link it to an existing archive policy

	createOpts := metrics.CreateOpts{
		ArchivePolicyName: "low",
		CreatedByProjectID: "3d40ca37-7234-4911-8987b9f288f4ae84",
		CreatedByUserID: "fdcfb420-c096-45e6-9e177a0bb1950884",
		Creator: "fdcfb420-c096-45e6-9e177a0bb1950884:3d40ca37-7234-4911-8987b9f288f4ae84",
	}
	metric, err := metrics.Create(gnocchiClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example of Creating a metric without an archive policy, assuming that Gnocchi has the needed
archive policy rule and can assign the policy automatically

	createOpts := metrics.CreateOpts{
		Unit: "MB",
	}
	metric, err := metrics.Create(gnocchiClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
*/
package metrics
