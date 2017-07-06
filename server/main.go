package main

import (
	"flag"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/volmex/volmex"
	"os"
	"os/signal"
	"syscall"
)

const pluginSockDir = "/run/docker/plugins"
const pluginSockName = "volmex"

func main() {
	storage := flag.String("storage", "/var/local/volmex", "Base for volume storage directories")
	flag.Parse()

	// check if volume storage base exists
	if _, err := os.Stat(*storage); os.IsNotExist(err) {
		panic(err)
	}

	// check if another instance of volmex is running
	if f, _ := os.Stat(pluginSockDir + "/" + pluginSockName + ".sock"); f != nil {
		panic("Plugin socket exists, is volmex running already?")
	}

	// initialize volmex volume state and the driver
	state := volmex.NewFileState(*storage + "/volumes.json")
	d := volmex.NewDriver(state, *storage)
	h := volume.NewHandler(d)

	// catch SIGINT / SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// delete plugin socket
			err := os.Remove(pluginSockDir + "/" + pluginSockName + ".sock")
			if err != nil {
				panic(err)
			}

			// terminate accordingly to the signal
			switch sig {
			case syscall.SIGINT:
				os.Exit(1)
			case syscall.SIGTERM:
				os.Exit(0)
			}
		}
	}()

	// serve driver as docker plugin socket
	err := h.ServeUnix(pluginSockName, 0)
	if err != nil {
		panic(err)
	}
}
