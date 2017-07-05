package main

import (
	"flag"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/volmex/volmex"
)

func main() {
	mountSource := flag.String("config", "/var/local/volmex", "Base for config directories")
	flag.Parse()

	config := volmex.NewFileVolConfig(*mountSource + "/volumes.json")
	d := volmex.NewDriver(config, *mountSource)
	h := volume.NewHandler(d)
	err := h.ServeUnix("volmex", 0)
	if err != nil {
		panic(err)
	}
}
