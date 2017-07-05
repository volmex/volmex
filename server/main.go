package main

import (
	"flag"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/volmex/volmex"
)

func main() {
	mountSource := flag.String("mountsource", "/var/local/volmex", "The base directory for the mounts")
	flag.Parse()

	storage := volmex.NewFileStorage(*mountSource + "/volumes.json")
	d := volmex.NewDriver(storage, *mountSource)
	h := volume.NewHandler(d)
	err := h.ServeUnix("volmex", 0)
	if err != nil {
		panic(err)
	}
}
