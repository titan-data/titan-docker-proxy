package proxy

import (
	"fmt"
	titan "github.com/titan-data/titan-client-go"
	"net/http"
)

/*
 * The proxy class is responsible for taking docker requests as input, and making the appropriate calls to an
 * instance of titan-server. The inputs to these functions are all structures defined in this package. The
 * responsibility of listening on the appropriate docker socket, marshalling to and from JSON, etc rests with
 * other portions of the package.
 */

type proxy struct {
	client     *titan.APIClient
	volumeApi  *titan.VolumesApiService
}

/*
 * /VolumeDriver.Capabilities
 *
 * This always returns a static definition with a "local" scope.
 */
func (p proxy) VolumeDriverCapabilities() VolumeDriverCapabilities {
	return VolumeDriverCapabilities{Capabilities:Capability{Scope:"local"}}
}

/*
 * /Plugin.Activate
 *
 * This always returns a static definition implementing "VolumeDriver"
 */
func (p proxy) PluginActivate() PluginDescription {
	return PluginDescription {
		Implements: []string{"VolumeDriver"},
	}
}

/*
 * Public proxy constructor. Takes a host ("localhost") and port (5001) to pass to the client.
 */
func Proxy(host string, port int) proxy {
	config := titan.NewConfiguration()
	config.Host = fmt.Sprintf("%s:%d", host, port)
	client := titan.NewAPIClient(config)
	return proxy{
		client:    client,
		volumeApi: client.VolumesApi,
	}
}

/*
 * For use in testing, this allows the test to pass a (mock) HTTP client to the titan client in order to facilitate
 * testing.
 */
func MockProxy(httpClient *http.Client) proxy {
	config := titan.NewConfiguration()
	config.HTTPClient = httpClient
	client := titan.NewAPIClient(config)
	return proxy{
		client:    client,
		volumeApi: client.VolumesApi,
	}
}