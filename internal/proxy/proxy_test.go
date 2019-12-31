package proxy

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

func TestPluginActivate(t *testing.T) {
	p := Proxy("localhost", 5001)
	desc := p.PluginActivate()
	assert.Equal(t, desc.Implements[0], "VolumeDriver")
}

func TestVolumeDriverCapabilities(t *testing.T) {
	p := Proxy("localhost", 5001)
	capabilities := p.VolumeDriverCapabilities()
	assert.Equal(t, capabilities.Capabilities.Scope, "local")
}

func TestListEmptyVolumes(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.RequestURI, "/v1/repositories/foo/volumes")
		w.Write([]byte("[]"))
	})
	client, teardown := testClient(h)
	defer teardown()

	resp, _ := client.Get("http://somehost/v1/repositories/foo/volumes")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), "[]")
}