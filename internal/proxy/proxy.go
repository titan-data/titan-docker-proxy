package proxy

import (
	"fmt"
	titan "github.com/titan-data/titan-client-go"
)

type proxy struct {
	client     *titan.APIClient
	volumeApi  *titan.VolumesApiService
}

func (p proxy) VolumeDriverCapabilities() VolumeDriverCapabilities {
	return VolumeDriverCapabilities{Capabilities:Capability{Scope:"local"}}
}

func (p proxy) PluginActivate() PluginDescription {
	return PluginDescription {
		Implements: []string{"VolumeDriver"},
	}
}

func Proxy(host string, port int) proxy {
	config := titan.NewConfiguration()
	config.Host = fmt.Sprintf("%s:%d", host, port)
	client := titan.NewAPIClient(config)
	return proxy{
		client:    client,
		volumeApi: client.VolumesApi,
	}
}

