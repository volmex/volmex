package main

import (
	"flag"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/volmex/volmex"
)

func main() {
	mountSource := flag.String("storage", "/var/local/volmex", "Base for storage directories")
	flag.Parse()

	state := volmex.NewFileState(*mountSource + "/volumes.json")
	d := volmex.NewDriver(state, *mountSource)
	h := volume.NewHandler(d)
	err := h.ServeUnix("volmex", 0)
	if err != nil {
		panic(err)
	}
}
