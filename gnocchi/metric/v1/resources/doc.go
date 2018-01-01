/*
Package resources provides the ability to retrieve resources through the Gnocchi API.

Example of Listing resources

	listOpts := resources.ListOpts{
		Details: True,
		ResourceType: "instance",
	}

	allPages, err := resources.List(gnocchiClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allResources, err := resources.ExtractResources(allPages)
	if err != nil {
		panic(err)
	}

	for _, resource := range allResources {
		fmt.Printf("%+v\n", resource)
	}
*/
package resources
