package testing

import (
	"testing"

	"github.com/gophercloud/utils/openstack/clientconfig"

	th "github.com/gophercloud/gophercloud/testhelper"
	yaml "gopkg.in/yaml.v2"
)

var VirginiaExpected = `clouds:
  virginia:
    auth:
      auth_url: https://va.example.com:5000/v3
      application_credential_id: app-cred-id
      application_credential_secret: secret
    auth_type: v3applicationcredential
    region_name: VA
    verify: true
`

var HawaiiExpected = `clouds:
  hawaii:
    auth:
      auth_url: https://hi.example.com:5000/v3
      username: jdoe
      password: password
      project_name: Some Project
      domain_name: default
    region_name: HNL
    verify: true
`

func TestMarshallCloudToYaml(t *testing.T) {
	clouds := make(map[string]map[string]*clientconfig.Cloud)
	clouds["clouds"] = map[string]*clientconfig.Cloud{
		"virginia": &VirginiaCloudYAML,
	}

	marshalled, err := yaml.Marshal(clouds)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, VirginiaExpected, string(marshalled))

	clouds["clouds"] = map[string]*clientconfig.Cloud{
		"hawaii": &HawaiiCloudYAML,
	}

	marshalled, err = yaml.Marshal(clouds)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, HawaiiExpected, string(marshalled))
}
