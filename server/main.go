package main

import (
	"flag"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/volmex/volmex"
	"os"
)

func main() {
	storage := flag.String("storage", "/var/local/volmex", "Base for storage directories")
	flag.Parse()

	_, err := os.Stat(*storage)
	if os.IsNotExist(err) {
		panic(err)
	}

	state := volmex.NewFileState(*storage + "/volumes.json")
	d := volmex.NewDriver(state, *storage)
	h := volume.NewHandler(d)
	err = h.ServeUnix("volmex", 0)
	if err != nil {
		panic(err)
	}
}
