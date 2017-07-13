package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/volmex/volmex"
)

const (
	pluginSockDir  = "/run/docker/plugins"
	pluginSockName = "volmex"
	pluginSock     = pluginSockDir + "/" + pluginSockName + ".sock"
)

func main() {
	storage := flag.String("storage", "/var/local/volmex", "base for volume storage directories")
	flag.Parse()

	// check if volume storage base exists
	if _, err := os.Stat(*storage); os.IsNotExist(err) {
		log.Fatalf("storage folder does not exist: %v", err)
	}

	// check if another instance of volmex is running
	if f, _ := os.Stat(pluginSock); f != nil {
		log.Fatalf("%v already exits, is volmex already running?", pluginSock)
	}

	// catch SIGINT / SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// delete plugin socket
			err := os.Remove(pluginSock)
			if err != nil {
				log.Printf("could not delete plugin socket: %v", err)
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

	// initialize volmex volume state and the driver
	state := volmex.NewFileState(*storage + "/volumes.json")
	d := volmex.NewDriver(state, *storage)
	h := volume.NewHandler(d)

	// serve driver as docker plugin socket
	err := h.ServeUnix(pluginSockName, 0)
	if err != nil {
		panic(err)
	}
}
