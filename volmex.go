package volmex

import (
	"fmt"

	"github.com/docker/go-plugins-helpers/volume"
	"os/exec"
	"strings"
)

type Volume struct {
	volume.Volume
	Options map[string]string
}

type Driver struct {
	storage     Storage
	mountSource string
}

func NewDriver(storage Storage, mountSource string) *Driver {
	return &Driver{
		storage:     storage,
		mountSource: mountSource,
	}
}

func (d *Driver) Create(req volume.Request) volume.Response {
	fmt.Printf("Create with %v\n", req)

	if req.Options["cmd"] == "" {
		return volume.Response{
			Err: "no mount command. Specify with -o cmd=acommand",
		}
	}

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

	fmt.Println("executing " + v.Options["cmd"])
	components := strings.Split(v.Options["cmd"], " ")
	cmd := exec.Command(components[0], components[1:]...)
	cmd.Env = []string{"MOUNT_SOURCE=" + v.Mountpoint}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	fmt.Println(string(out))

	return volume.Response{
		Mountpoint: v.Mountpoint,
	}
}

func (d *Driver) Unmount(req volume.UnmountRequest) volume.Response {
	fmt.Printf("Unmount with %v\n", req)
	return volume.Response{Err: "no such volume"}
}

func (d *Driver) Capabilities(req volume.Request) volume.Response {
	fmt.Printf("Capabilities with %v\n", req)
	return volume.Response{}
}
