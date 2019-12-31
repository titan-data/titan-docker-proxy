package proxy

import (
	"testing"
	assert "github.com/stretchr/testify/assert"
)

var p = Proxy("localhost", 5001)

func TestPluginActivate(t *testing.T) {
	assert := assert.New(t)

	desc := p.PluginActivate()
	assert.Equal(desc.Implements[0], "VolumeDriver")
}

func TestVolumeDriverCapabilities(t *testing.T) {
	assert := assert.New(t)

	cap := p.VolumeDriverCapabilities()
	assert.Equal(cap.Capabilities.Scope, "local")
}