package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

// Logger is an interface representing the Logger struct
type Logger interface {
	Printf(format string, args ...interface{})
}

// DefaultLogger is a default struct, which satisfies the Logger interface
type DefaultLogger struct{}

// Printf is a default Printf method
func (DefaultLogger) Printf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

// RoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper
type RoundTripper struct {
	// Default http.RoundTripper
	Rt http.RoundTripper
	// If set to true, the HTTP "Cache-Control: no-cache" header will be added by
	// default to all API requests
	NoCacheHeader bool
	// How many times HTTP connection should be retried until giving up
	MaxRetries int
	// If Logger is not nil, then RoundTrip method will debug the JSON requests
	// and responses
	Logger Logger
}

// List of headers that contain sensitive data.
var sensitiveHeaders = map[string]struct{}{
	"x-auth-token":                    {},
	"x-auth-key":                      {},
	"x-service-token":                 {},
	"x-storage-token":                 {},
	"x-account-meta-temp-url-key":     {},
	"x-account-meta-temp-url-key-2":   {},
	"x-container-meta-temp-url-key":   {},
	"x-container-meta-temp-url-key-2": {},
	"set-cookie":                      {},
	"x-subject-token":                 {},
}

func hideSensitiveHeadersData(headers http.Header) []string {
	result := make([]string, len(headers))
	headerIdx := 0
	for header, data := range headers {
		if _, ok := sensitiveHeaders[strings.ToLower(header)]; ok {
			result[headerIdx] = fmt.Sprintf("%s: %s", header, "***")
		} else {
			result[headerIdx] = fmt.Sprintf("%s: %s", header, strings.Join(data, " "))
		}
		headerIdx++
	}

	return result
}

// formatHeaders converts standard http.Header type to a string with separated headers.
// It will hide data of sensitive headers.
func formatHeaders(headers http.Header, separator string) string {
	redactedHeaders := hideSensitiveHeadersData(headers)
	sort.Strings(redactedHeaders)

	return strings.Join(redactedHeaders, separator)
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	defer func() {
		if request.Body != nil {
			request.Body.Close()
		}
	}()

	// for future reference, this is how to access the Transport struct:
	//tlsconfig := lrt.Rt.(*http.Transport).TLSClientConfig

	if lrt.NoCacheHeader {
		// set Cache-Control header to no-cache on request to force HTTP caches (if any) to go upstream.
		// This is a work-around until all Openstack APIs implement proper Cache-Control headers by their own.
		// The guidelines for this were added to http://specs.openstack.org/openstack/api-sig/guidelines/http/caching.html in 03/2018.
		request.Header.Set("Cache-Control", "no-cache")
	}

	var err error

	if lrt.Logger != nil {
		lrt.Logger.Printf("OpenStack Request URL: %s %s", request.Method, request.URL)
		lrt.Logger.Printf("OpenStack Request Headers:\n%s", formatHeaders(request.Header, "\n"))

		if request.Body != nil {
			request.Body, err = lrt.logRequest(request.Body, request.Header.Get("Content-Type"))
			if err != nil {
				return nil, err
			}
		}
	}

	response, err := lrt.Rt.RoundTrip(request)

	// If the first request didn't return a response, retry up to `max_retries`.
	retry := 1
	for response == nil {
		if retry > lrt.MaxRetries {
			if lrt.Logger != nil {
				lrt.Logger.Printf("OpenStack connection error, retries exhausted. Aborting")
			}
			err = fmt.Errorf("OpenStack connection error, retries exhausted. Aborting. Last error was: %s", err)
			return nil, err
		}

		if lrt.Logger != nil {
			lrt.Logger.Printf("OpenStack connection error, retry number %d: %s", retry, err)
		}
		response, err = lrt.Rt.RoundTrip(request)
		retry += 1
	}

	if lrt.Logger != nil {
		lrt.Logger.Printf("OpenStack Response Code: %d", response.StatusCode)
		lrt.Logger.Printf("OpenStack Response Headers:\n%s", formatHeaders(response.Header, "\n"))

		response.Body, err = lrt.logResponse(response.Body, response.Header.Get("Content-Type"))
	}

	return response, err
}

// logRequest will log the HTTP Request details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *RoundTripper) logRequest(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	// Handle request contentType
	if strings.HasPrefix(contentType, "application/json") {
		var bs bytes.Buffer
		defer original.Close()

		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}

		debugInfo, err := FormatJSON(bs.Bytes())
		if err != nil {
			lrt.Logger.Printf("%s", err)
		}
		lrt.Logger.Printf("OpenStack Request Body: %s", debugInfo)

		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}

	lrt.Logger.Printf("Not logging because OpenStack request body isn't JSON")
	return original, nil
}

// logResponse will log the HTTP Response details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (lrt *RoundTripper) logResponse(original io.ReadCloser, contentType string) (io.ReadCloser, error) {
	if strings.HasPrefix(contentType, "application/json") {
		var bs bytes.Buffer
		defer original.Close()

		_, err := io.Copy(&bs, original)
		if err != nil {
			return nil, err
		}

		debugInfo, err := FormatJSON(bs.Bytes())
		if err != nil {
			lrt.Logger.Printf("%s", err)
		}
		if debugInfo != "" {
			lrt.Logger.Printf("OpenStack Response Body: %s", debugInfo)
		}

		return ioutil.NopCloser(strings.NewReader(bs.String())), nil
	}

	lrt.Logger.Printf("Not logging because OpenStack response body isn't JSON")
	return original, nil
}

// FormatJSON will try to pretty-format a JSON body.
// It will also mask known fields which contain sensitive information.
func FormatJSON(raw []byte) (string, error) {
	var rawData interface{}

	err := json.Unmarshal(raw, &rawData)
	if err != nil {
		return string(raw), fmt.Errorf("unable to parse OpenStack JSON: %s", err)
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		pretty, err := json.MarshalIndent(rawData, "", "  ")
		if err != nil {
			return string(raw), fmt.Errorf("unable to re-marshal OpenStack JSON: %s", err)
		}

		return string(pretty), nil
	}

	// Mask known password fields
	if v, ok := data["auth"].(map[string]interface{}); ok {
		// v2 auth methods
		if v, ok := v["passwordCredentials"].(map[string]interface{}); ok {
			v["password"] = "***"
		}
		if v, ok := v["token"].(map[string]interface{}); ok {
			v["id"] = "***"
		}
		// v3 auth methods
		if v, ok := v["identity"].(map[string]interface{}); ok {
			if v, ok := v["password"].(map[string]interface{}); ok {
				if v, ok := v["user"].(map[string]interface{}); ok {
					v["password"] = "***"
				}
			}
			if v, ok := v["application_credential"].(map[string]interface{}); ok {
				v["secret"] = "***"
			}
			if v, ok := v["token"].(map[string]interface{}); ok {
				v["id"] = "***"
			}
		}
	}

	// Ignore the huge catalog output
	if v, ok := data["token"].(map[string]interface{}); ok {
		if _, ok := v["catalog"]; ok {
			v["catalog"] = "***"
		}
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return string(raw), fmt.Errorf("unable to re-marshal OpenStack JSON: %s", err)
	}

	return string(pretty), nil
}
