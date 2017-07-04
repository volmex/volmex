package volmex

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-plugins-helpers/volume"
)

const (
	manifest         = `{"Implements": ["VolumeDriver"]}`
	createPath       = "/VolumeDriver.Create"
	getPath          = "/VolumeDriver.Get"
	listPath         = "/VolumeDriver.List"
	removePath       = "/VolumeDriver.Remove"
	hostVirtualPath  = "/VolumeDriver.Path"
	mountPath        = "/VolumeDriver.Mount"
	unmountPath      = "/VolumeDriver.Unmount"
	capabilitiesPath = "/VolumeDriver.Capabilities"
)

func TestEmpty(t *testing.T) {
	d := &VolmexDriver{}
	h := volume.NewHandler(d)

	l := sockets.NewInmemSocket("test", 0)
	go h.Serve(l)
	defer l.Close()

	client := &http.Client{Transport: &http.Transport{
		Dial: l.Dial,
	}}

	// Create
	resp, err := driverRequest(client, createPath, volume.Request{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err == "" {
		t.Fatal("Did not throw an error")
	}
}

// Initiates a new request to the driver
func driverRequest(client *http.Client, method string, req volume.Request) (*volume.Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp volume.Response
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}
