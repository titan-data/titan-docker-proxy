package proxy

import (
	"context"
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
	ctx        context.Context
}

/*
 * Converts an error object into an "Err" string to return to consumers. If this is a titan-server API error, then
 * we return the message field. Otherwise, we return the default error string.
 */
func getErrorString(err error) string  {
	if openApiErr, ok := err.(titan.GenericOpenAPIError); ok {
		if apiErr, ok := openApiErr.Model().(titan.ApiError); ok {
			return apiErr.Message
		}
	}
	return err.Error()
}

/*
 * Converts from a Titan volume to a Docker volume. The main difference is that the repository name is part of the
 * volume name. The mountpoint is also pulled out of the properties to a first class response.
 */
func titanToDocker(repo string, vol titan.Volume) Volume {
	return Volume {
		Name: fmt.Sprintf("%s/%s", repo, vol.Name),
		Mountpoint: vol.Config["mountpoint"].(string),
		Status: map[string]string{},
	}
}

/*
 * /VolumeDriver.Capabilities
 *
 * This always returns a static definition with a "local" scope.
 */
func (p proxy) VolumeCapabilities() VolumeCapabilities {
	return VolumeCapabilities{Capabilities: Capability{Scope: "local"}}
}

/*
 * /VolumeDriver.List
 *
 * Returns a list of all volumes on the system. This requires iterating over all repositories followed by the volumes
 * for each.
 */
func (p proxy) ListVolumes() ListVolumeResponse {
	repositories, _, err := p.client.RepositoriesApi.ListRepositories(p.ctx)
	if err != nil {
		return ListVolumeResponse{Err: getErrorString(err)}
	}

	ret := ListVolumeResponse{
		Volumes: []Volume{},
	}

	for _, repo := range repositories {
		volumes, _, err := p.client.VolumesApi.ListVolumes(p.ctx, repo.Name)
		if err != nil {
			return ListVolumeResponse{Err: getErrorString(err)}
		}
		for _, vol := range volumes {
			ret.Volumes = append(ret.Volumes, titanToDocker(repo.Name, vol))
		}
	}

	return ret
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
		client:     client,
		ctx:		context.Background(),
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
		ctx:       context.Background(),
	}
}