package volmex

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/go-plugins-helpers/volume"
)

// Errors
var (
	NoMountCommandErr = errors.New("no volmex mount command. specify with -o cmd=\"command\"")
	UnknownVolumeErr  = errors.New("volume unknown")
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

// the following methods implement the docker volume plugin protocol (v1)
// https://docs.docker.com/engine/extend/plugins_volume/#volumedriver

// Create is called by docker upon 'docker volume create' and creates a new volmex.Volume
func (d *Driver) Create(req *volume.CreateRequest) error {
	fmt.Printf("received 'Create' request for volume: %v\n", req.Name)

	// load driver state
	err := d.state.Load()
	if err != nil {
		return err
	}

	// check if a command was given
	if req.Options["cmd"] == "" {
		return NoMountCommandErr
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
		return err
	}

	fmt.Printf("\tmountpoint: %v\n", v.Volume.Mountpoint)
	fmt.Printf("\tcommand: %v\n", v.Options["cmd"])

	return nil
}

// Get returns the volume configuration for a given volume name
func (d *Driver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	fmt.Printf("received 'Get' request for volume: %v\n", req.Name)

	// load driver state
	err := d.state.Load()
	if err != nil {
		return &volume.GetResponse{}, err
	}

	// try to retrieve volume configuration from driver state
	v, err := d.state.Get(req.Name)
	if err != nil {
		return &volume.GetResponse{}, err
	}

	return &volume.GetResponse{
		Volume: &v.Volume,
	}, nil
}

// List returns all known volume configurations
func (d *Driver) List() (*volume.ListResponse, error) {
	fmt.Printf("received 'List' request\n")

	// load driver state
	err := d.state.Load()
	if err != nil {
		return &volume.ListResponse{}, err
	}

	// retrieve volume configurations from state
	var vs []*volume.Volume
	for _, v := range d.state.List() {
		vs = append(vs, &v.Volume)
	}

	fmt.Printf("\t%v known volumes\n", len(vs))
	for _, v := range vs {
		fmt.Printf("\t\t%v\n", v.Name)
	}

	return &volume.ListResponse{
		Volumes: vs,
	}, nil
}

// Remove is issued by docker when a user requests to delete a volume
// however, we don't actually remove any files or the storage folder
func (d *Driver) Remove(req *volume.RemoveRequest) error {
	fmt.Printf("received 'Remove' request for volume: %v\n", req.Name)

	// load driver state
	err := d.state.Load()
	if err != nil {
		return err
	}

	// remove volume from state
	d.state.Remove(req.Name)

	// save driver state
	err = d.state.Save()
	if err != nil {
		return err
	}

	return nil
}

// Path returns the mountpoint for a given volume name
func (d *Driver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	fmt.Printf("received 'Path' request for volume: %v\n", req.Name)

	// load driver state
	err := d.state.Load()
	if err != nil {
		return &volume.PathResponse{}, err
	}

	// try to retrieve volume configuration from driver state
	v, err := d.state.Get(req.Name)
	if err != nil {
		return &volume.PathResponse{}, UnknownVolumeErr
	}

	return &volume.PathResponse{
		Mountpoint: v.Mountpoint,
	}, nil
}

// Mount is called by docker before a container using a volmex volume is started
// since we don't actually do anything related to storage, we only execute the specified volmex command and return the mountpoint
func (d *Driver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	fmt.Printf("received 'Mount' request for volume: %v\n", req.Name)

	// load driver state
	err := d.state.Load()
	if err != nil {
		return &volume.MountResponse{}, err
	}

	// try to retrieve volume configuration from driver state
	v, err := d.state.Get(req.Name)
	if err != nil {
		return &volume.MountResponse{}, UnknownVolumeErr
	}

	fmt.Println("\texecuting volmex command: " + v.Options["cmd"])

	// prepare command
	cmdString := strings.TrimSpace(v.Options["cmd"])
	cmdParts := strings.Split(cmdString, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	// set volmex environment for the new command
	cmd.Env = []string{
		"VOLMEX_NAME=" + v.Name,
		"VOLMEX_MOUNTPOINT=" + v.Mountpoint,
		"VOLMEX_CMD=" + cmdString,
	}

	// execute command and print output
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error from volmex command: %v", err.Error())
	}
	fmt.Println(string(out))

	return &volume.MountResponse{
		Mountpoint: v.Mountpoint,
	}, nil
}

// Unmount is called by docker after a container using a volmex volume was stopped
// Currently we don't want to execute a command after a container was stopped, so we don't need to do anything here
func (d *Driver) Unmount(req *volume.UnmountRequest) error {
	fmt.Printf("received 'Unmount' request for volume: %v\n", req.Name)
	return nil
}

// Capabilities is called by docker to get certain driver options (atm only the 'scope')
func (d *Driver) Capabilities() *volume.CapabilitiesResponse {
	fmt.Printf("received 'Capabilities' request")
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "local"}}
}
