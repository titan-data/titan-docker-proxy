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
 * JSON, and invoking the appropriate methods of the forwarder to then make calls to titan-server. Because this is so
 * generic, we are able to implement a generic interposition layer and use reflection to do all the work.
 */

type Listener interface {
	Listen() error
}

type listener struct {
	forw forwarder.Forwarder
	path string
	mux  *http.ServeMux
}

type handler struct {
	listen *listener
	req    interface{}
	fun    interface{}
}

/*
 * The main handler method. This will detect whether the method expects zero arguments or one, handles marshaling,
 * and errors while invoking the given method on the forwarder.
 */
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var response []reflect.Value
	var err error

	funcValue := reflect.ValueOf(h.fun)
	if h.req != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			err = json.Unmarshal(body, h.req)
		}

		response = funcValue.Call([]reflect.Value{reflect.ValueOf(h.req).Elem()})
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

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func create(forward forwarder.Forwarder, path string) listener {
	l := &listener{
		forw: forward,
		path: path,
		mux:  http.NewServeMux(),
	}

	l.mux.Handle("/VolumeDriver.Create", handler{l, &forwarder.CreateVolumeRequest{}, forward.CreateVolume})
	l.mux.Handle("/VolumeDriver.Get", handler{l, &forwarder.VolumeRequest{}, forward.GetVolume})
	l.mux.Handle("/VolumeDriver.Path", handler{l, &forwarder.VolumeRequest{}, forward.GetPath})
	l.mux.Handle("/VolumeDriver.List", handler{l, nil, forward.ListVolumes})
	l.mux.Handle("/VolumeDriver.Mount", handler{l, &forwarder.MountVolumeRequest{}, forward.MountVolume})
	l.mux.Handle("/VolumeDriver.Remove", handler{l, &forwarder.VolumeRequest{}, forward.RemoveVolume})
	l.mux.Handle("/VolumeDriver.Unmount", handler{l, &forwarder.MountVolumeRequest{}, forward.UnmountVolume})

	return *l
}

func (l listener) Listen() error {
	listen, err := net.Listen("unix", l.path)
	if err != nil {
		return fmt.Errorf("listen failed on %s: %w", l.path, err)
	}

	return http.Serve(listen, l.mux)
}

func New(forwarder forwarder.Forwarder, path string) Listener {
	return create(forwarder, path)
}
