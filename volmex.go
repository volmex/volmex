package volmex

import (
	"errors"
	"fmt"

	"github.com/docker/go-plugins-helpers/volume"
)

type VolmexDriver struct {
	volumes     []*volume.Volume
	mountSource string
}

func New(mountSource string) *VolmexDriver {
	return &VolmexDriver{
		mountSource: mountSource,
	}
}

func (d *VolmexDriver) Create(req volume.Request) volume.Response {
	fmt.Printf("Create with %v\n", req)
	d.volumes = append(d.volumes,
		&volume.Volume{
			Name:       req.Name,
			Mountpoint: d.mountSource + "/" + req.Name,
		})
	return volume.Response{}
}

func (d *VolmexDriver) Get(req volume.Request) volume.Response {
	fmt.Printf("Get with %v", req)

	v, err := volumeByName(d.volumes, req.Name)
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	return volume.Response{
		Volume: v,
	}
}

func (d *VolmexDriver) List(req volume.Request) volume.Response {
	fmt.Printf("List with %v\n", req)
	return volume.Response{
		Volumes: d.volumes,
	}
}

func (d *VolmexDriver) Remove(req volume.Request) volume.Response {
	fmt.Printf("Remove with %v\n", req)
	for i := range d.volumes {
		if d.volumes[i].Name == req.Name {
			d.volumes = append(d.volumes[:i], d.volumes[i+1:]...)
			return volume.Response{}
		}
	}
	return volume.Response{}
}

func (d *VolmexDriver) Path(req volume.Request) volume.Response {
	fmt.Printf("Path with %v\n", req)

	v, err := volumeByName(d.volumes, req.Name)
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	return volume.Response{
		Mountpoint: v.Mountpoint,
	}
}

func (d *VolmexDriver) Mount(req volume.MountRequest) volume.Response {
	fmt.Printf("Mount with %v\n", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Unmount(req volume.UnmountRequest) volume.Response {
	fmt.Printf("Unmount with %v\n", req)
	return volume.Response{Err: "no such volume"}
}

func (d *VolmexDriver) Capabilities(req volume.Request) volume.Response {
	fmt.Printf("Capabilities with %v\n", req)
	return volume.Response{}
}

func volumeByName(volumes []*volume.Volume, name string) (*volume.Volume, error) {
	for _, v := range volumes {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, errors.New("no volume found")
}
