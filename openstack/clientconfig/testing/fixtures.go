package testing

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

var HawaiiCloudYAML = clientconfig.Cloud{
	RegionName: "HNL",
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://hi.example.com:5000/v3",
		Username:    "jdoe",
		Password:    "password",
		ProjectName: "Some Project",
		DomainName:  "default",
	},
}

var HawaiiClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://hi.example.com:5000/v3",
		Username:    "jdoe",
		Password:    "password",
		ProjectName: "Some Project",
		DomainName:  "default",
	},
}

var HawaiiEnvAuth = map[string]string{
	"OS_AUTH_URL":     "https://hi.example.com:5000/v3",
	"OS_USERNAME":     "jdoe",
	"OS_PASSWORD":     "password",
	"OS_PROJECT_NAME": "Some Project",
	"OS_DOMAIN_NAME":  "default",
}

var HawaiiAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectName: "Some Project",
		DomainName:  "default",
	},
	IdentityEndpoint: "https://hi.example.com:5000/v3",
	Username:         "jdoe",
	Password:         "password",
	TenantName:       "Some Project",
	DomainName:       "default",
}

var FloridaCloudYAML = clientconfig.Cloud{
	RegionName: "MIA",
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:      "https://fl.example.com:5000/v3",
		Username:     "jdoe",
		Password:     "password",
		ProjectID:    "12345",
		UserDomainID: "abcde",
	},
}

var FloridaClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:      "https://fl.example.com:5000/v3",
		Username:     "jdoe",
		Password:     "password",
		ProjectID:    "12345",
		UserDomainID: "abcde",
	},
}

var FloridaEnvAuth = map[string]string{
	"OS_AUTH_URL":       "https://fl.example.com:5000/v3",
	"OS_USERNAME":       "jdoe",
	"OS_PASSWORD":       "password",
	"OS_PROJECT_ID":     "12345",
	"OS_USER_DOMAIN_ID": "abcde",
}

var FloridaAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectID: "12345",
	},
	IdentityEndpoint: "https://fl.example.com:5000/v3",
	Username:         "jdoe",
	Password:         "password",
	TenantID:         "12345",
	DomainID:         "abcde",
}

var CaliforniaCloudYAML = clientconfig.Cloud{
	Regions: []interface{}{
		"SAN",
		"LAX",
	},
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:           "https://ca.example.com:5000/v3",
		Username:          "jdoe",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "default",
		UserDomainName:    "default",
	},
}

var CaliforniaClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:           "https://ca.example.com:5000/v3",
		Username:          "jdoe",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "default",
		UserDomainName:    "default",
	},
}

var CaliforniaEnvAuth = map[string]string{
	"OS_AUTH_URL":            "https://ca.example.com:5000/v3",
	"OS_USERNAME":            "jdoe",
	"OS_PASSWORD":            "password",
	"OS_PROJECT_NAME":        "Some Project",
	"OS_PROJECT_DOMAIN_NAME": "default",
	"OS_USER_DOMAIN_NAME":    "default",
}

var CaliforniaAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectName: "Some Project",
		DomainName:  "default",
	},
	IdentityEndpoint: "https://ca.example.com:5000/v3",
	Username:         "jdoe",
	Password:         "password",
	TenantName:       "Some Project",
	DomainName:       "default",
}

var ArizonaCloudYAML = clientconfig.Cloud{
	RegionName: "PHX",
	AuthType:   clientconfig.AuthToken,
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://az.example.com:5000/v3",
		Token:       "12345",
		ProjectName: "Some Project",
		DomainName:  "default",
	},
}

var ArizonaClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://az.example.com:5000/v3",
		Token:       "12345",
		ProjectName: "Some Project",
		DomainName:  "default",
	},
}

var ArizonaEnvAuth = map[string]string{
	"OS_AUTH_URL":     "https://az.example.com:5000/v3",
	"OS_TOKEN":        "12345",
	"OS_PROJECT_NAME": "Some Project",
	"OS_DOMAIN_NAME":  "default",
}

var ArizonaAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectName: "Some Project",
		DomainName:  "default",
	},
	IdentityEndpoint: "https://az.example.com:5000/v3",
	TokenID:          "12345",
	TenantName:       "Some Project",
}

var NewMexicoCloudYAML = clientconfig.Cloud{
	RegionName: "SAF",
	AuthType:   clientconfig.AuthPassword,
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:           "https://nm.example.com:5000/v3",
		Username:          "jdoe",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "Some Domain",
		UserDomainName:    "Some OtherDomain",
		DomainName:        "default",
	},
}

var NewMexicoClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:           "https://nm.example.com:5000/v3",
		Username:          "jdoe",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "Some Domain",
		UserDomainName:    "Other Domain",
		DomainName:        "default",
	},
}

var NewMexicoEnvAuth = map[string]string{
	"OS_AUTH_URL":            "https://nm.example.com:5000/v3",
	"OS_USERNAME":            "jdoe",
	"OS_PASSWORD":            "password",
	"OS_PROJECT_NAME":        "Some Project",
	"OS_PROJECT_DOMAIN_NAME": "Some Domain",
	"OS_USER_DOMAIN_NAME":    "Other Domain",
	"OS_DOMAIN_NAME":         "default",
}

var NewMexicoAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectName: "Some Project",
		DomainName:  "Some Domain",
	},
	IdentityEndpoint: "https://nm.example.com:5000/v3",
	Username:         "jdoe",
	Password:         "password",
	TenantName:       "Some Project",
	DomainName:       "Other Domain",
}

var NevadaCloudYAML = clientconfig.Cloud{
	RegionName: "LAS",
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:           "https://nv.example.com:5000/v3",
		UserID:            "12345",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "Some Domain",
	},
}

var NevadaClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:           "https://nv.example.com:5000/v3",
		UserID:            "12345",
		Password:          "password",
		ProjectName:       "Some Project",
		ProjectDomainName: "Some Domain",
	},
}

var NevadaEnvAuth = map[string]string{
	"OS_AUTH_URL":            "https://nv.example.com:5000/v3",
	"OS_USER_ID":             "12345",
	"OS_PASSWORD":            "password",
	"OS_PROJECT_NAME":        "Some Project",
	"OS_PROJECT_DOMAIN_NAME": "Some Domain",
}

var NevadaAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectName: "Some Project",
		DomainName:  "Some Domain",
	},
	IdentityEndpoint: "https://nv.example.com:5000/v3",
	UserID:           "12345",
	Password:         "password",
	TenantName:       "Some Project",
}

var TexasCloudYAML = clientconfig.Cloud{
	RegionName: "AUS",
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:        "https://tx.example.com:5000/v3",
		Username:       "jdoe",
		Password:       "password",
		ProjectName:    "Some Project",
		UserDomainName: "Some Domain",
		DefaultDomain:  "default",
	},
}

var TexasClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:        "https://tx.example.com:5000/v3",
		Username:       "jdoe",
		Password:       "password",
		ProjectName:    "Some Project",
		UserDomainName: "Some Domain",
		DefaultDomain:  "default",
	},
}

var TexasEnvAuth = map[string]string{
	"OS_AUTH_URL":         "https://tx.example.com:5000/v3",
	"OS_USERNAME":         "jdoe",
	"OS_PASSWORD":         "password",
	"OS_PROJECT_NAME":     "Some Project",
	"OS_USER_DOMAIN_NAME": "Some Domain",
	"OS_DEFAULT_DOMAIN":   "default",
}

var TexasAuthOpts = &gophercloud.AuthOptions{
	Scope: &gophercloud.AuthScope{
		ProjectName: "Some Project",
		DomainID:    "default",
	},
	IdentityEndpoint: "https://tx.example.com:5000/v3",
	Username:         "jdoe",
	Password:         "password",
	TenantName:       "Some Project",
	DomainName:       "Some Domain",
}

var CloudYAML = clientconfig.Clouds{
	Clouds: map[string]clientconfig.Cloud{
		"hawaii":     HawaiiCloudYAML,
		"florida":    FloridaCloudYAML,
		"california": CaliforniaCloudYAML,
		"arizona":    ArizonaCloudYAML,
		"newmexico":  NewMexicoCloudYAML,
		"nevada":     NevadaCloudYAML,
		"texas":      TexasCloudYAML,
	},
}

var AlbertaCloudYAML = clientconfig.Cloud{
	RegionName: "YYC",
	AuthType:   clientconfig.AuthPassword,
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://ab.example.com:5000/v2.0",
		Username:    "jdoe",
		Password:    "password",
		ProjectName: "Some Project",
	},
}

var AlbertaClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://ab.example.com:5000/v2.0",
		Username:    "jdoe",
		Password:    "password",
		ProjectName: "Some Project",
	},
}

var AlbertaEnvAuth = map[string]string{
	"OS_AUTH_URL":             "https://ab.example.com:5000/v2.0",
	"OS_USERNAME":             "jdoe",
	"OS_PASSWORD":             "password",
	"OS_PROJECT_NAME":         "Some Project",
	"OS_IDENTITY_API_VERSION": "2.0",
}

var AlbertaAuthOpts = &gophercloud.AuthOptions{
	IdentityEndpoint: "https://ab.example.com:5000/v2.0",
	Username:         "jdoe",
	Password:         "password",
	TenantName:       "Some Project",
}

var YukonCloudYAML = clientconfig.Cloud{
	RegionName: "YXY",
	AuthType:   clientconfig.AuthV2Token,
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://yt.example.com:5000/v2.0",
		Token:       "12345",
		ProjectName: "Some Project",
	},
}

var YukonClientOpts = &clientconfig.ClientOpts{
	AuthInfo: &clientconfig.AuthInfo{
		AuthURL:     "https://yt.example.com:5000/v2.0",
		Token:       "12345",
		ProjectName: "Some Project",
	},
}

var YukonEnvAuth = map[string]string{
	"OS_AUTH_URL":             "https://yt.example.com:5000/v2.0",
	"OS_TOKEN":                "12345",
	"OS_PROJECT_NAME":         "Some Project",
	"OS_IDENTITY_API_VERSION": "2",
}

var YukonAuthOpts = &gophercloud.AuthOptions{
	IdentityEndpoint: "https://yt.example.com:5000/v2.0",
	TokenID:          "12345",
	TenantName:       "Some Project",
}

var LegacyCloudYAML = clientconfig.Clouds{
	Clouds: map[string]clientconfig.Cloud{
		"alberta": AlbertaCloudYAML,
		"yukon":   YukonCloudYAML,
	},
}
