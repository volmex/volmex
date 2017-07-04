package volmex

import (
	"errors"
	"fmt"

	"github.com/docker/go-plugins-helpers/volume"
)

type Volume struct {
	volume.Volume
	Options map[string]string
}

type Driver struct {
	storage     Storage
	mountSource string
}

func (d *Driver) Create(req volume.Request) volume.Response {
	fmt.Printf("Create with %v\n", req)
	v := &Volume{
		Volume: volume.Volume{
			Name:       req.Name,
			Mountpoint: d.mountSource + "/" + req.Name,
		},
		Options: req.Options,
	}
	d.storage.Put(v.Name, v)
	return volume.Response{}
}

func (d *Driver) Get(req volume.Request) volume.Response {
	fmt.Printf("Get with %v", req)

	v, err := d.storage.Get(req.Name)
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	return volume.Response{
		Volume: &v.Volume,
	}
}

func (d *Driver) List(req volume.Request) volume.Response {
	fmt.Printf("List with %v\n", req)
	var vs []*volume.Volume
	for _, v := range d.storage.List() {
		vs = append(vs, &v.Volume)
	}
	return volume.Response{
		Volumes: vs,
	}
}

func (d *Driver) Remove(req volume.Request) volume.Response {
	fmt.Printf("Remove with %v\n", req)
	d.storage.Remove(req.Name)
	return volume.Response{}
}

func (d *Driver) Path(req volume.Request) volume.Response {
	fmt.Printf("Path with %v\n", req)

	v, err := d.storage.Get(req.Name)
	if err != nil {
		return volume.Response{
			Err: "no volume found",
		}
	}

	return volume.Response{
		Mountpoint: v.Mountpoint,
	}
}

func (d *Driver) Mount(req volume.MountRequest) volume.Response {
	fmt.Printf("Mount with %v\n", req)

	v, err := d.storage.Get(req.Name)
	if err != nil {
		return volume.Response{
			Err: "no volume found",
		}
	}

	if v.Options["cmd"] == "" {
		return volume.Response{
			Err: "no mount command. Specify with -o cmd=acommand",
		}
	}
	return volume.Response{}
}

func (d *Driver) Unmount(req volume.UnmountRequest) volume.Response {
	fmt.Printf("Unmount with %v\n", req)
	return volume.Response{Err: "no such volume"}
}

func (d *Driver) Capabilities(req volume.Request) volume.Response {
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
