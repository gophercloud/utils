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

Example of Creating measures inside different metrics via metric ID references in one request

	currentTimestamp := time.Now().UTC()
	pastHourTimestamp := currentTimestamp.Add(-1 * time.Hour)
	createOpts := measures.BatchCreateMetricsOpts{
		{
			ID: "777a01d6-4694-49cb-b86a-5ba9fd4e609e",
			Measures: []measures.MeasureOpts{
				{
					Timestamp: &currentTimestamp,
					Value:     200,
				},
				{
					Timestamp: &pastHourTimestamp,
					Value:     300,
				},
			},
		},
		{
			ID: "6dbc97c5-bfdf-47a2-b184-02e7fa348d21",
			Measures: []measures.MeasureOpts{
				{
					Timestamp: &currentTimestamp,
					Value:     111,
				},
				{
					Timestamp: &pastHourTimestamp,
					Value:     222,
				},
			},
		},
	}
	if err := measures.BatchCreateMetrics(gnocchiClient, createOpts).ExtractErr(); err != nil && err.Error() != "EOF" {
		panic(err)
	}

Example of Creating measures inside different metrics via metric names and resource IDreferences of that metrics in one request

	currentTimestamp := time.Now().UTC()
	pastHourTimestamp := currentTimestamp.Add(-1 * time.Hour)
	createOpts := measures.CreateBatchResourcesMetricsOpts{
		BatchResourcesMetricsMeasuresOpts: map[string]map[string][]measures.MeasureOpts{
			"1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc": {
				"memory.usage": []measures.MeasureOpts{
					{
						Timestamp: &currentTimestamp,
						Value:     1562.82,
					},
					{
						Timestamp: &pastHourTimestamp,
						Value:     768.1,
					},
				},
			},
			"789a7f65-977d-40f4-beed-f717100125f5": {
				"cpu.util": []measures.MeasureOpts{
					{
						Timestamp: &currentTimestamp,
						Value:     89.9,
					},
					{
						Timestamp: &pastHourTimestamp,
						Value:     56,
					},
				},
				"network.incoming.bytes.rate": []measures.MeasureOpts{
					{
						Timestamp: &currentTimestamp,
						Value:     15671.32,
					},
					{
						Timestamp: &pastHourTimestamp,
						Value:     87123,
					},
				},
			},
		},
	}
	if err := measures.CreateBatchResourcesMetrics(gnocchiClient, createOpts).ExtractErr(); err != nil && err.Error() != "EOF" {
		panic(err)
	}
*/
package measures
