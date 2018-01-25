/*
Package measures provides the ability to retrieve measures through the Gnocchi API.

Example of Listing measures of a known metric

	startTime := time.Date(2018, 1, 4, 10, 0, 0, 0, time.UTC)
	metricID := "9e5a6441-1044-4181-b66e-34e180753040"
	listOpts := measures.ListOpts{
		Resample: "2h",
		Granularity: "1h",
		Start: &startTime,
	}
	allPages, err := measures.List(gnocchiClient, metricID, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allMeasures, err := measures.ExtractMeasures(allPages)
	if err != nil {
		panic(err)
	}

	for _, measure := range allMeasures {
		fmt.Printf("%+v\n", measure)
	}

Example of Creating measures inside a single metric

	createOpts := measures.CreateOpts{
		Measures: []measures.MeasureOpts{
			{
				Timestamp: time.Date(2018, 1, 18, 12, 31, 0, 0, time.UTC),
				Value:     101.2,
			},
			{
				Timestamp: time.Date(2018, 1, 18, 14, 32, 0, 0, time.UTC),
				Value:     102,
			},
		},
	}
	metricID := "9e5a6441-1044-4181-b66e-34e180753040"
	if err := measures.Create(gnocchiClient, metricID, createOpts).ExtractErr(); err != nil && err.Error() != "EOF" {
		panic(err)
	}
*/
package measures
