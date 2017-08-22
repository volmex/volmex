package volmex

import (
	"bytes"
	"encoding/json"
	"io"
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
	resp, err := driverRequest(client, createPath, volume.CreateRequest{Name: "foo", Options: map[string]string{"cmd": "foo"}})
	if err != nil {
		t.Fatal(err)
	}
	var errResp *volume.ErrorResponse
	if err := json.NewDecoder(resp).Decode(&errResp); err != nil {
		t.Fatal(err)
	}
	if errResp.Err != "" {
		t.Fatal(err)
	}

	// Should not create with missing opts: cmd
	resp, err = driverRequest(client, createPath, volume.CreateRequest{Name: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if err = json.NewDecoder(resp).Decode(&errResp); err != nil {
		t.Fatal(err)
	}
	if errResp.Err != NoMountCommandErr.Error() {
		t.Fatal(errResp.Err)
	}

	// Test Get
	resp, err = driverRequest(client, getPath, volume.GetRequest{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	var getResp *volume.GetResponse
	if err := json.NewDecoder(resp).Decode(&getResp); err != nil {
		t.Fatal(err)
	}
	if getResp.Volume.Mountpoint != "/var/local/volmex/foo" {
		t.Fatalf("resp.Volume.Mountpoint was wrong %v", getResp.Volume.Mountpoint)
	}

	// Test List
	resp, err = driverRequest(client, listPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	var listResp *volume.ListResponse
	if err := json.NewDecoder(resp).Decode(&listResp); err != nil {
		t.Fatal(err)
	}
	if listResp.Volumes[0].Name != "foo" {
		t.Fatalf("List did not contain volume %v", listResp.Volumes)
	}

	// Test Path
	resp, err = driverRequest(client, hostVirtualPath, volume.PathRequest{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	var pathResp *volume.PathResponse
	if err := json.NewDecoder(resp).Decode(&pathResp); err != nil {
		t.Fatal(err)
	}
	if pathResp.Mountpoint != "/var/local/volmex/foo" {
		t.Fatalf("resp.Mountpoint was not %v", pathResp.Mountpoint)
	}

	// Test Remove
	resp, err = driverRequest(client, removePath, volume.RemoveRequest{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(resp).Decode(&errResp); err != nil {
		t.Fatal(err)
	}

	// Get removed volume should fail
	resp, err = driverRequest(client, getPath, volume.GetRequest{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(resp).Decode(&getResp); err != nil {
		t.Fatal(err)
	}
}

// initiates a new request to the driver
func driverRequest(client *http.Client, method string, req interface{}) (io.Reader, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	if req == nil {
		b = []byte{}
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
