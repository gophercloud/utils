/*
Package resourcetypes provides ability to manage resource types through the Gnocchi API.

Example of Listing resource types

  allPages, err := resourcetypes.List(client).AllPages()
  if err != nil {
    panic(err)
	}

  allResourceTypes, err := resourcetypes.ExtractResourceTypes(allPages)
  if err != nil {
    panic(err)
	}

  for _, resourceType := range allResourceTypes {
    fmt.Printf("%+v\n", resourceType)
  }
*/
package resourcetypes
