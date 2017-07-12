package volmex

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/go-plugins-helpers/volume"
)

// Volume holds volmex related volume configurations
// - Volume from go-plugins-helpers/volume
// - Options from 'docker volume create' (volume.Request from Create)
type Volume struct {
	volume.Volume
	Options map[string]string
}

// Driver represents the driver's configuration
// - State holds the drivers's volume configurations
// - mountBase is the base folder for volume storage directories
type Driver struct {
	state     State
	mountBase string
}

// NewDriver creates a new Driver instance
func NewDriver(state State, mountBase string) *Driver {
	return &Driver{
		state:     state,
		mountBase: mountBase,
	}
}

// Create is called by docker upon 'docker volume create' and creates a new volmex.Volume
func (d *Driver) Create(req volume.Request) volume.Response {
	fmt.Printf("Create Request for volume: %v\n", req.Name)

	// load driver state
	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	// check if a command was given
	if req.Options["cmd"] == "" {
		return volume.Response{
			Err: "no mount command. specify with -o cmd=\"command\"",
		}
	}

	// check if the volume storage folder is missing and create it if necessary
	if _, err := os.Stat(d.mountBase + "/" + req.Name); os.IsNotExist(err) {
		os.Mkdir(d.mountBase+"/"+req.Name, 0777)
	}

	// store the volume configuration in a volmex.Volume and add it to the driver's state
	v := &Volume{
		Volume: volume.Volume{
			Name:       req.Name,
			Mountpoint: d.mountBase + "/" + req.Name,
		},
		Options: req.Options,
	}
	d.state.Put(v.Name, v)

	// save driver state
	err = d.state.Save()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	fmt.Printf("\tmountpoint: %v\n", v.Volume.Mountpoint)
	fmt.Printf("\tcommand: %v\n", v.Options["cmd"])

	return volume.Response{}
}

func (d *Driver) Get(req volume.Request) volume.Response {
	fmt.Printf("Get with %v", req)

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	v, err := d.state.Get(req.Name)
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

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	var vs []*volume.Volume
	for _, v := range d.state.List() {
		vs = append(vs, &v.Volume)
	}
	return volume.Response{
		Volumes: vs,
	}
}

func (d *Driver) Remove(req volume.Request) volume.Response {
	fmt.Printf("Remove with %v\n", req)

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	d.state.Remove(req.Name)

	err = d.state.Save()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	return volume.Response{}
}

func (d *Driver) Path(req volume.Request) volume.Response {
	fmt.Printf("Path with %v\n", req)

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	v, err := d.state.Get(req.Name)
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

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	v, err := d.state.Get(req.Name)
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

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	return volume.Response{Err: "no such volume"}
}

func (d *Driver) Capabilities(req volume.Request) volume.Response {
	fmt.Printf("Capabilities with %v\n", req)

	err := d.state.Load()
	if err != nil {
		return volume.Response{
			Err: err.Error(),
		}
	}

	return volume.Response{}
}
