/*
Package resources provides the ability to retrieve resources through the Gnocchi API.

Example of Listing resources

	resourceType: "instance",
	listOpts := resources.ListOpts{
		Details: True,
	}

	allPages, err := resources.List(gnocchiClient, listOpts, resourceType).AllPages()
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
