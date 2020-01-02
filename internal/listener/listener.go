/*
 * Copyright The Titan Project Contributors.
 */
package listener

import (
	"encoding/json"
	"errors"
	"github.com/titan-data/titan-docker-proxy/internal/forwarder"
)

/*
 * The listener is responsible for listening on a Unix Domain Socket for docker requests, marshaling data to and from
 * JSON, and invoking the appropriate methods of the forwarder to then make calls to titan-server.
 */

type Listener interface {
	Listen() error
}

type listener struct {
	Forwarder forwarder.Forwarder
	Path      string
}

func (l listener) CreateVolume(body []byte) []byte {
	var request forwarder.CreateVolumeRequest
	var responseBody []byte

	err := json.Unmarshal(body, &request)
	if err == nil {
		response := l.Forwarder.CreateVolume(request)
		responseBody, err = json.Marshal(response)
	}

	if err != nil {
		responseBody, _ = json.Marshal(forwarder.VolumeResponse{Err: err.Error()})
	}
	return responseBody
}

func (l listener) Listen() error {
	return errors.New("TODO")
}

func Raw(forwarder forwarder.Forwarder, path string) listener {
	return listener{
		Forwarder: forwarder,
		Path:      path,
	}
}

func New(forwarder forwarder.Forwarder, path string) Listener {
	return Raw(forwarder, path)
}
