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
	desc := p.PluginActivate()
	assert.Equal(t, desc.Implements[0], "VolumeDriver")
}

func TestVolumeDriverCapabilities(t *testing.T) {
	p := Proxy("localhost", 5001)
	capabilities := p.VolumeCapabilities()
	assert.Equal(t, capabilities.Capabilities.Scope, "local")
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

	volumes := p.ListVolumes()
	if (assert.Empty(t, volumes.Err) &&
		assert.Equal(t, len(volumes.Volumes), 2)) {
		assert.Equal(t, volumes.Volumes[0].Name, "foo/v0")
		assert.Equal(t, volumes.Volumes[0].Mountpoint, "/v0")
		assert.Equal(t, len(volumes.Volumes[0].Status), 0)
		assert.Equal(t, volumes.Volumes[1].Name, "foo/v1")
		assert.Equal(t, volumes.Volumes[1].Mountpoint, "/v1")
		assert.Equal(t, len(volumes.Volumes[1].Status), 0)
	}
}