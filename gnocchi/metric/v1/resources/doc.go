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
	resourceType = "generic"
	resource, err := resources.Get(gnocchiClient, resourceType, resourceID).Extract()
	if err != nil {
		panic(err)
	}

Example of Creating a resource without a metric

	createOpts := resources.CreateOpts{
		ID: "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55",
		ProjectID: "4154f088-8333-4e04-94c4-1155c33c0fc9",
		UserID: "bd5874d666624b24a9f01c128871e4ac",
	}
	resourceType = ""
	resource, err := resources.Create(gnocchiClient, createOpts, resourceType).Extract()
	if err != nil {
		panic(err)
	}

Example of Creating a resource with links to some existing metrics

	createOpts := resources.CreateOpts{
		ID: "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55",
		ProjectID: "4154f088-8333-4e04-94c4-1155c33c0fc9",
		UserID: "bd5874d666624b24a9f01c128871e4ac",
		Metrics: map[string]string{
			"disk.read.bytes.rate": "ed1bb76f-6ccc-4ad2-994c-dbb19ddccbae",
			"disk.write.bytes.rate": "0a2da84d-4753-43f5-a65f-0f8d44d2766c",
		},
	}
	resourceType = "compute_instance_disk"
	resource, err := resources.Create(gnocchiClient, createOpts, resourceType).Extract()
	if err != nil {
		panic(err)
	}

Example of Creating a resource and a metric a the same time

	createOpts := resources.CreateOpts{
		ID: "23d5d3f7-9dfa-4f73-b72b-8b0b0063ec55",
		ProjectID: "4154f088-8333-4e04-94c4-1155c33c0fc9",
		UserID: "bd5874d666624b24a9f01c128871e4ac",
		Metrics: map[string]map[string]string{
			"cpu.delta": map[string]string{
				"archive_policy_name": "medium",
			},
		},
	}
	resourceType = "compute_instance"
	resource, err := resources.Create(gnocchiClient, createOpts, resourceType).Extract()
	if err != nil {
		panic(err)
	}
*/
package resources
