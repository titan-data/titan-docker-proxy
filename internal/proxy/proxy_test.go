package proxy

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testProxy(handler http.Handler) (proxy, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return MockProxy(cli), s.Close
}

func TestPluginActivate(t *testing.T) {
	p := Proxy("localhost", 5001)
	resp := p.PluginActivate()
	assert.Equal(t, resp.Implements[0], "VolumeDriver")
}

func TestVolumeDriverCapabilities(t *testing.T) {
	p := Proxy("localhost", 5001)
	resp := p.VolumeCapabilities()
	assert.Equal(t, resp.Capabilities.Scope, "local")
}

func TestListVolumes(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/v1/repositories" {
			w.Write([]byte("[{\"name\":\"foo\",\"properties\":{}}]"))
		} else {
			assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes")
			w.Write([]byte("[{\"name\":\"v0\",\"config\":{\"mountpoint\":\"/v0\"}}," +
				"{\"name\":\"v1\",\"config\":{\"mountpoint\":\"/v1\"}}]"))
		}
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.ListVolumes()
	if assert.Empty(t, resp.Err) &&
		assert.Equal(t, len(resp.Volumes), 2) {
		assert.Equal(t, resp.Volumes[0].Name, "foo/v0")
		assert.Equal(t, resp.Volumes[0].Mountpoint, "/v0")
		assert.Equal(t, len(resp.Volumes[0].Status), 0)
		assert.Equal(t, resp.Volumes[1].Name, "foo/v1")
		assert.Equal(t, resp.Volumes[1].Mountpoint, "/v1")
		assert.Equal(t, len(resp.Volumes[1].Status), 0)
	}
}

func TestListVolumesRepoError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, r.RequestURI, "/v1/repositories")
		w.WriteHeader(404)
		w.Write([]byte("{\"message\":\"no such repository\"}"))
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.ListVolumes()
	assert.Equal(t, resp.Err, "no such repository")
}

func TestListVolumesVolumeError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.RequestURI == "/v1/repositories" {
			w.Write([]byte("[{\"name\":\"foo\",\"properties\":{}}]"))
		} else {
			assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes")
			w.WriteHeader(404)
			w.Write([]byte("{\"message\":\"no such volume\"}"))
		}
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.ListVolumes()
	assert.Equal(t, resp.Err, "no such volume")
}

func TestGetVolume(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes/vol")
		w.Write([]byte("{\"name\":\"vol\",\"config\":{\"mountpoint\":\"/vol\"}}"))
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.GetVolume(VolumeRequest{Name: "foo/vol"})
	if assert.Empty(t, resp.Err) {
		assert.Equal(t, resp.Volume.Name, "foo/vol")
		assert.Equal(t, resp.Volume.Mountpoint, "/vol")
	}
}

func TestGetVolumeBadName(t *testing.T) {
	p := Proxy("localhost", 5001)

	resp := p.GetVolume(VolumeRequest{Name: "foo"})
	assert.Equal(t, resp.Err, "volume name must be of the form <repository>/<volume>")
}

func TestGetVolumeError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes/vol")
		w.WriteHeader(404)
		w.Write([]byte("{\"message\":\"no such volume\"}"))
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.GetVolume(VolumeRequest{Name: "foo/vol"})
	assert.Equal(t, resp.Err, "no such volume")
}

func TestGetPath(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes/vol")
		w.Write([]byte("{\"name\":\"vol\",\"config\":{\"mountpoint\":\"/vol\"}}"))
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.GetPath(VolumeRequest{Name: "foo/vol"})
	if assert.Empty(t, resp.Err) {
		assert.Equal(t, resp.Mountpoint, "/vol")
	}
}

func TestGetPathError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes/vol")
		w.WriteHeader(404)
		w.Write([]byte("{\"message\":\"no such volume\"}"))
	})
	p, teardown := testProxy(h)
	defer teardown()

	resp := p.GetPath(VolumeRequest{Name: "foo/vol"})
	assert.Equal(t, resp.Err, "no such volume")
}
