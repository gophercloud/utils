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

Example of Getting a resource

	resourceID = "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55"
	resourceType = ""
	resource, err := resources.Get(gnocchiClient, resourceID, resourceType).Extract()
	if err != nil {
		panic(err)
	}

*/
package resources
