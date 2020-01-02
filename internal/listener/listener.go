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
	"reflect"
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

type handler struct {
	l         *listener
	r         interface{}
	f         interface{}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var response []reflect.Value
	var err error

	funcValue := reflect.ValueOf(h.f)

	funcType := reflect.TypeOf(h.f)
	argCount := funcType.NumIn()
	if (argCount == 1) {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			err = json.Unmarshal(body, h.r)
		}

		response = funcValue.Call([]reflect.Value{ reflect.ValueOf(h.r).Elem() })
	} else {
		response = funcValue.Call([]reflect.Value{})
	}

	var body []byte
	if err == nil {
		response = reflect.ValueOf(json.Marshal).Call(response)
		body = response[0].Bytes()
		if !response[1].IsNil() {
			err = response[1].Interface().(error)
		}
	}

	if err != nil {
		body, err = json.Marshal(forwarder.VolumeResponse{Err: err.Error()})
	}

	if body == nil {
		body = []byte("{\"Err\":\"Unable to serialize error response\"}")
	}

	w.Write(body)
}

func create(forward forwarder.Forwarder, path string) listener {
	l := &listener{
		Forwarder: forward,
		Path:      path,
		Mux:       http.NewServeMux(),
	}
	l.Mux.Handle("/VolumeDriver.Create", handler{l, &forwarder.CreateVolumeRequest{}, forward.CreateVolume})
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
