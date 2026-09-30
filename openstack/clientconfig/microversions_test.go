package clientconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultMicroversionSerialization(t *testing.T) {
	const input = `auth:
  auth_url: https://example.org/v3
compute_default_microversion: "2.87"
volumev3_default_microversion: "3.60"
block_storage_default_microversion: null
shared_file_system_default_microversion: 2.10
future_service_default_microversion: "1.4"
unrelated_setting:
  nested: true
regions:
  - name: mars
    values:
      compute_default_microversion: "2.79"
`
	var cloud Cloud
	if err := yaml.Unmarshal([]byte(input), &cloud); err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"compute": "2.87", "volumev3": "3.60", "shared-file-system": "2.10", "future-service": "1.4"}
	if !reflect.DeepEqual(cloud.DefaultMicroversions, expected) {
		t.Fatalf("defaults: %#v", cloud.DefaultMicroversions)
	}
	for _, codec := range []struct {
		name      string
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{{"yaml", yaml.Marshal, yaml.Unmarshal}, {"json", json.Marshal, json.Unmarshal}} {
		t.Run(codec.name, func(t *testing.T) {
			data, err := codec.marshal(cloud)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "volumev3_default_microversion") {
				t.Fatalf("missing flat key: %s", data)
			}
			var got Cloud
			if err := codec.unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cloud, got) {
				t.Fatalf("roundtrip: %#v", got)
			}
		})
	}
	for _, input := range []string{"compute_default_microversion: [2.87]", "compute_default_microversion: {bad: value}"} {
		if err := yaml.Unmarshal([]byte(input), &Cloud{}); err == nil {
			t.Fatalf("accepted invalid value: %s", input)
		}
	}
}

type microversionYAML struct{ cloud, secure, public Cloud }

func (m microversionYAML) LoadCloudsYAML() (map[string]Cloud, error) {
	return map[string]Cloud{"test": m.cloud}, nil
}
func (m microversionYAML) LoadSecureCloudsYAML() (map[string]Cloud, error) {
	return map[string]Cloud{"test": m.secure}, nil
}
func (m microversionYAML) LoadPublicCloudsYAML() (map[string]Cloud, error) {
	return map[string]Cloud{"example": m.public}, nil
}

func TestConfiguredMicroversionsReachRequests(t *testing.T) {
	t.Setenv("OS_CLOUD", "")
	for _, tc := range []struct{ service, serviceType, key, version, header, value string }{
		{"compute", "compute", "compute", "2.79", "OpenStack-API-Version", "compute 2.79"},
		{"volume", "volumev3", "volumev3", "3.60", "X-OpenStack-Volume-API-Version", "3.60"},
		{"sharev2", "sharev2", "shared-file-system", "2.65", "X-OpenStack-Manila-API-Version", "2.65"},
		{"placement", "placement", "placement", "1.20", "OpenStack-API-Version", "placement 1.20"},
		{"baremetal", "baremetal", "bare-metal", "1.80", "X-OpenStack-Ironic-API-Version", "1.80"},
		{"baremetal-introspection", "baremetal-introspection", "baremetal-introspection", "1.12", "X-OpenStack-Ironic-Inspector-API-Version", "1.12"},
		{"container", "container", "application-container", "1.40", "OpenStack-API-Version", "application-container 1.40"},
		{"container-infra", "container-infra", "container-infrastructure-management", "1.4", "OpenStack-API-Version", "container-infra 1.4"},
		{"image", "image", "compute", "", "OpenStack-API-Version", ""},
	} {
		t.Run(tc.service, func(t *testing.T) {
			for _, explicit := range []string{"", "9.9"} {
				t.Run("override="+explicit, func(t *testing.T) {
					mux := http.NewServeMux()
					server := httptest.NewServer(mux)
					defer server.Close()
					major := "2"
					if tc.version != "" {
						major = strings.SplitN(tc.version, ".", 2)[0]
					}
					mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Subject-Token", "test-token")
						w.WriteHeader(http.StatusCreated)
						fmt.Fprintf(w, `{"token":{"expires_at":"2100-01-01T00:00:00Z","catalog":[{"type":%q,"name":"service","endpoints":[{"interface":"public","region":"mars","url":%q}]}]}}`, tc.serviceType, server.URL+"/service/v"+major+"/")
					})
					called := false
					mux.HandleFunc("/service/", func(w http.ResponseWriter, r *http.Request) {
						called = true
						expected := tc.value
						if explicit != "" {
							if tc.header == "OpenStack-API-Version" {
								prefix := tc.serviceType
								if tc.value != "" {
									prefix = strings.SplitN(tc.value, " ", 2)[0]
								}
								expected = prefix + " " + explicit
							} else {
								expected = explicit
							}
						}
						if got := r.Header.Get(tc.header); got != expected {
							t.Errorf("header: %q, want %q", got, expected)
						}
						w.WriteHeader(http.StatusOK)
					})
					var cloud Cloud
					if err := yaml.Unmarshal([]byte("compute_default_microversion: \"2.87\"\n"), &cloud); err != nil {
						t.Fatal(err)
					}
					cloud.AuthInfo = &AuthInfo{AuthURL: server.URL + "/v3", Username: "user", Password: "password", DomainName: "default", ProjectName: "project"}
					cloud.RegionName = "mars"
					if tc.service != "image" {
						cloud.DefaultMicroversions[tc.key] = tc.version
					}
					if tc.service == "compute" {
						cloud.DefaultMicroversions["compute"] = "2.87"
						cloud.Regions = []Region{{Name: "mars", Values: Cloud{DefaultMicroversions: map[string]string{"compute": tc.version}}}}
					}
					client, err := NewServiceClient(context.Background(), tc.service, &ClientOpts{Cloud: "test", RegionName: "mars", Microversion: explicit, YAMLOpts: microversionYAML{cloud: cloud}})
					if err != nil {
						t.Fatal(err)
					}
					expected := tc.version
					if explicit != "" {
						expected = explicit
					}
					if client.Microversion != expected {
						t.Fatalf("microversion: %q, want %q", client.Microversion, expected)
					}
					if _, err := client.Get(context.Background(), server.URL+"/service/", nil, nil); err != nil {
						t.Fatal(err)
					}
					if !called {
						t.Fatal("service request not received")
					}
				})
			}
		})
	}
}

// Expectations are generated by the SDK itself; Python is not required for Go tests.
func TestOpenStackSDKDefaultMicroversions(t *testing.T) {
	data, err := os.ReadFile("testdata/openstacksdk-defaults.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Defaults struct {
			Region, Service string
			Cloud           map[string]any
		}
		Cases []struct {
			Name                      string
			Cloud, Secure, Public     map[string]any
			Region, Service, Expected string
		}
	}
	if err := yaml.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	for _, tc := range fixture.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Region == "" {
				tc.Region = fixture.Defaults.Region
			}
			if tc.Service == "" {
				tc.Service = fixture.Defaults.Service
			}
			fields := maps.Clone(fixture.Defaults.Cloud)
			maps.Copy(fields, tc.Cloud)
			tc.Cloud = fields

			var input microversionYAML
			for _, pair := range []struct {
				data  map[string]any
				cloud *Cloud
			}{{tc.Cloud, &input.cloud}, {tc.Secure, &input.secure}, {tc.Public, &input.public}} {
				data, err := yaml.Marshal(pair.data)
				if err != nil {
					t.Fatal(err)
				}
				if err := yaml.Unmarshal(data, pair.cloud); err != nil {
					t.Fatal(err)
				}
			}
			cloud, err := GetCloudFromYAML(&ClientOpts{Cloud: "test", RegionName: tc.Region, YAMLOpts: input})
			if err != nil {
				t.Fatal(err)
			}
			got := cloud.DefaultMicroversion(tc.Service)
			if got != tc.Expected {
				t.Fatalf("SDK expects %q, got %q", tc.Expected, got)
			}
		})
	}
}
