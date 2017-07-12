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

func TestVolmex(t *testing.T) {
	state := NewInMemoryState()
	d := &Driver{
		state:     state,
		mountBase: "/var/local/volmex",
	}
	h := volume.NewHandler(d)

	l := sockets.NewInmemSocket("", 0)
	go h.Serve(l)
	defer l.Close()

	client := &http.Client{Transport: &http.Transport{
		Dial: l.Dial,
	}}

	// Test Create
	resp, err := driverRequest(client, createPath, volume.Request{Name: "foo", Options: map[string]string{"cmd": "foo"}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatal(resp.Err)
	}

	// Test Create with missing command
	resp, err = driverRequest(client, createPath, volume.Request{Name: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err == "" {
		t.Fatal("Create without cmd worked")
	}

	// Test Get
	resp, err = driverRequest(client, getPath, volume.Request{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatal(resp.Err)
	}
	if resp.Volume.Mountpoint != "/var/local/volmex/foo" {
		t.Fatalf("resp.Volume.Mountpoint was wrong %v", resp.Volume.Mountpoint)
	}

	// Test List
	resp, err = driverRequest(client, listPath, volume.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatal(resp.Err)
	}
	if resp.Volumes[0].Name != "foo" {
		t.Fatalf("List did not contain volume %v", resp.Volumes)
	}

	// Test Path
	resp, err = driverRequest(client, hostVirtualPath, volume.Request{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatal(resp.Err)
	}
	if resp.Mountpoint != "/var/local/volmex/foo" {
		t.Fatalf("resp.Mountpoint was not %v", resp.Mountpoint)
	}

	// Test Remove
	resp, err = driverRequest(client, removePath, volume.Request{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatal(resp.Err)
	}

	// Test removed Get
	resp, err = driverRequest(client, getPath, volume.Request{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err == "" {
		t.Fatal("volume was not deleted")
	}
}

// Initiates a new request to the driver
func driverRequest(client *http.Client, path string, req volume.Request) (*volume.Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+path, "application/json", bytes.NewReader(b))
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
