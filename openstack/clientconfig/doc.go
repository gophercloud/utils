/*
Package clientconfig provides convienent functions for creating OpenStack
clients. It is based on the Python os-client-config library.

See https://docs.openstack.org/os-client-config/latest for details.

Example to Create a Provider Client From clouds.yaml

	opts := &clientconfig.ClientOpts{
		Cloud: "hawaii",
	}

	pClient, err := clientconfig.AuthenticatedClient(ctx, opts)
	if err != nil {
		panic(err)
	}

Example to Manually Create a Provider Client

	opts := &clientconfig.ClientOpts{
		AuthInfo: &clientconfig.AuthInfo{
			AuthURL:     "https://hi.example.com:5000/v3",
			Username:    "jdoe",
			Password:    "password",
			ProjectName: "Some Project",
			DomainName:  "default",
		},
	}

	pClient, err := clientconfig.AuthenticatedClient(ctx, opts)
	if err != nil {
		panic(err)
	}

Example to Create a Service Client from clouds.yaml

	opts := &clientconfig.ClientOpts{
		Cloud: "hawaii",
	}

	computeClient, err := clientconfig.NewServiceClient(ctx, "compute", opts)
	if err != nil {
		panic(err)
	}

# Service Default Microversions

NewServiceClient applies {service_type}_default_microversion settings from
clouds.yaml, for example compute_default_microversion: "2.87". Hyphens in service
types become underscores in configuration keys. Canonical types and aliases are
accepted; block_storage_default_microversion, block_store_default_microversion,
volume_default_microversion, volumev2_default_microversion and
volumev3_default_microversion all configure block storage. Canonical names take
precedence, followed by aliases in OpenStack Service Types Authority order.
For Cinder this is volumev3, volumev2, volume, then block-store. YAML null defaults
are skipped during lookup, while an explicitly empty string clears the default.

Defaults survive profile, secure.yaml and regional configuration merging.
ClientOpts.Microversion overrides the configured service default. An unset or
empty default leaves the microversion unset. Callers can also change
ServiceClient.Microversion after creation. Defaults are sent to the API without
automatic microversion negotiation.
*/
package clientconfig
