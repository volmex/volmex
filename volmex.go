package volmex

import (
	"fmt"

	"github.com/docker/go-plugins-helpers/volume"
)

type VolmexDriver struct {
	volumes []string
}

func (d *VolmexDriver) Create(req volume.Request) volume.Response {
	fmt.Printf("Create with %v", req)
	return volume.Response{Err: "no volume created"}
}

func (d *VolmexDriver) Get(req volume.Request) volume.Response {
	fmt.Printf("Get with %v", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) List(req volume.Request) volume.Response {
	fmt.Printf("List with %v", req)
	return volume.Response{Err: "no volumes"}
}

func (d *VolmexDriver) Remove(req volume.Request) volume.Response {
	fmt.Printf("Remove with %v", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Path(req volume.Request) volume.Response {
	fmt.Printf("Path with %v", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Mount(req volume.MountRequest) volume.Response {
	fmt.Printf("Mount with %v", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Unmount(req volume.UnmountRequest) volume.Response {
	fmt.Printf("Unmount with %v", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Capabilities(req volume.Request) volume.Response {
	fmt.Printf("Capabilities with %v", req)
	return volume.Response{}
}
