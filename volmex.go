package volmex

import "github.com/docker/go-plugins-helpers/volume"

type VolmexDriver struct {
	volumes []string
}

func (d *VolmexDriver) Create(req volume.Request) volume.Response {
	return volume.Response{Err: "no volume created"}
}

func (d *VolmexDriver) Get(req volume.Request) volume.Response {
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) List(req volume.Request) volume.Response {
	return volume.Response{Err: "no volumes"}
}

func (d *VolmexDriver) Remove(req volume.Request) volume.Response {
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Path(req volume.Request) volume.Response {
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Mount(req volume.MountRequest) volume.Response {
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Unmount(req volume.UnmountRequest) volume.Response {
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Capabilities(req volume.Request) volume.Response {
	return volume.Response{}
}
