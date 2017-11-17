package testing

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

var CloudYAMLHawaii = clientconfig.Cloud{
	RegionName: "HNL",
	Auth: &clientconfig.CloudAuth{
		AuthURL:     "https://hi.example.com:5000/v3",
		Username:    "jdoe",
		Password:    "password",
		ProjectName: "Some Project",
		DomainName:  "default",
	},
}

var CloudYAMLFlorida = clientconfig.Cloud{
	RegionName: "MIA",
	Auth: &clientconfig.CloudAuth{
		AuthURL:      "https://fl.example.com:5000/v3",
		Username:     "jdoe",
		Password:     "password",
		ProjectID:    "12345",
		UserDomainID: "abcde",
	},
}

var CloudYAMLCalifornia = clientconfig.Cloud{
	Regions: []interface{}{
		"SAN",
		"LAX",
	},
	Auth: &clientconfig.CloudAuth{
		AuthURL:           "https://ca.example.com:5000/v3",
		Username:          "jdoe",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "default",
	},
}

var CloudYAML = clientconfig.Clouds{
	Clouds: map[string]clientconfig.Cloud{
		"hawaii":     CloudYAMLHawaii,
		"florida":    CloudYAMLFlorida,
		"california": CloudYAMLCalifornia,
	},
}

var HawaiiAuthOpts = &gophercloud.AuthOptions{
	IdentityEndpoint: "https://hi.example.com:5000/v3",
	Username:         "jdoe",
	Password:         "password",
	TenantName:       "Some Project",
	DomainName:       "default",
}
