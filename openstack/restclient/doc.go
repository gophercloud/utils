/* Package restclient provides generic REST functions.

Example of a GET request

	getURL := computeClient.ServiceURL("flavors", flavorID)
	getRes := restclient.Get(computeClient, getURL, nil)
	if err != nil {
		panic(err)
	}

	flavorResult, err = getRes.Extract()
	if err != nil {
		panic(err)
	}

	flavor := flavorResult["flavor"].(map[string]interface{})
	fmt.Printf("%v\n", flavor)

Example of a POST request

	flavorOpts := map[string]interface{}{
		"name":  "some-name",
		"ram":   512,
		"vcpus": 1,
		"disk":  5,
		"id":    "some-id",
	}

	postOpts := &restclient.PostOpts{
		Params: map[string]interface{}{"flavor": flavorOpts},
	}

	postURL := computeClient.ServiceURL("flavors")
	postRes := restclient.Post(computeClient, postURL, postOpts)
	if err != nil {
		panic(err)
	}

	flavorResult, err := postRes.Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", flavor)

Example of a DELETE Request

	deleteURL := computeClient.ServiceURL("flavors", "flavor-id")
	deleteRes := restclient.Delete(computeClient, deleteURL, nil)
	if err != nil {
		panic(err)
	}

	err = deleteRes.ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package restclient
