/*
 * Copyright The Titan Project Contributors.
 */
package listener

import (
	"encoding/json"
	"fmt"
	"github.com/titan-data/titan-docker-proxy/internal/forwarder"
	"io/ioutil"
	"net"
	"net/http"
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
	Mux		  *http.ServeMux
}

func getBody(r *http.Request, request interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(body, request)
	}
	return err
}

func writeResponse(w http.ResponseWriter, response interface{}, err error) {
	var body []byte
	if err == nil {
		body, err = json.Marshal(response)
	}

	if err != nil {
		body, err = json.Marshal(forwarder.VolumeResponse{Err: err.Error()})
	}

	if body == nil {
		body = []byte("{\"Err\":\"Unable to serialize error response\"}")
	}

	w.Write(body)
}

func (l listener) CreateVolume(w http.ResponseWriter, r *http.Request) {
	var request forwarder.CreateVolumeRequest
	var response forwarder.VolumeResponse

	err := getBody(r, &request)
	if err == nil {
		response = l.Forwarder.CreateVolume(request)
	}

	writeResponse(w, response, err)
}

func create(forwarder forwarder.Forwarder, path string) listener {
	l := &listener{
		Forwarder: forwarder,
		Path:      path,
		Mux:       http.NewServeMux(),
	}
	l.Mux.HandleFunc("/VolumeDriver.Create", l.CreateVolume)
	return *l
}

func (l listener) Listen() error {
	listen, err := net.Listen("unix", l.Path)
	if err != nil {
		return fmt.Errorf("listen failed on %s: %w", l.Path, err)
	}

	return http.Serve(listen, l.Mux)
}

func New(forwarder forwarder.Forwarder, path string) Listener {
	return create(forwarder, path)
}
