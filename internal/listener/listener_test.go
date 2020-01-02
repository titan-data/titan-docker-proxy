/*
 * Copyright The Titan Project Contributors.
 */

package listener

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/titan-data/titan-docker-proxy/internal/forwarder"
	"testing"
)

type MockForwarder struct {
	mock.Mock
}

func (f *MockForwarder) CreateVolume(request forwarder.CreateVolumeRequest) forwarder.VolumeResponse {
	return forwarder.VolumeResponse{Err: ""}
}

func (f *MockForwarder) GetPath(request forwarder.VolumeRequest) forwarder.GetPathResponse {
	return forwarder.GetPathResponse{
		Mountpoint: "/" + request.Name,
		Err:        "",
	}
}

func (f *MockForwarder) GetVolume(request forwarder.VolumeRequest) forwarder.GetVolumeResponse {
	return forwarder.GetVolumeResponse{
		Err: "",
		Volume: forwarder.Volume{
			Name:       request.Name,
			Mountpoint: "/" + request.Name,
			Status:     map[string]string{},
		},
	}
}

func (f *MockForwarder) ListVolumes() forwarder.ListVolumeResponse {
	return forwarder.ListVolumeResponse{
		Err: "",
		Volumes: []forwarder.Volume{
			{
				Name:       "foo/vol",
				Mountpoint: "/foo/vol",
				Status:     map[string]string{},
			},
		},
	}
}

func (f *MockForwarder) MountVolume(request forwarder.MountVolumeRequest) forwarder.VolumeResponse {
	return forwarder.VolumeResponse{Err: ""}
}

func (f *MockForwarder) PluginActivate() forwarder.PluginDescription {
	return forwarder.PluginDescription{Implements: []string{"VolumeDriver"}}
}

func (f *MockForwarder) RemoveVolume(request forwarder.VolumeRequest) forwarder.VolumeResponse {
	return forwarder.VolumeResponse{Err: ""}
}

func (f *MockForwarder) VolumeCapabilities() forwarder.VolumeCapabilities {
	return forwarder.VolumeCapabilities{Capabilities: forwarder.Capability{Scope: "local"}}
}

func (f *MockForwarder) UnmountVolume(request forwarder.MountVolumeRequest) forwarder.VolumeResponse {
	return forwarder.VolumeResponse{Err: ""}
}

func TestCreateVolume(t *testing.T) {
	f := new(MockForwarder)
	l := Raw(f, "/path")
	resp := l.CreateVolume([]byte("{\"Name\":\"foo/vol\",\"Opts\":{}}"))
	assert.Equal(t, "{\"Err\":\"\"}", string(resp))
}
