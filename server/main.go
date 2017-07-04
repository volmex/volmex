package main

import (
	"github.com/docker/go-plugins-helpers/volume"
	"gitlab.mi.hdm-stuttgart.de/volmex"
)

func main() {
	d := &volmex.VolmexDriver{}
	h := volume.NewHandler(d)
	err := h.ServeUnix("volmex", 0)
	if err != nil {
		panic(err)
	}
}
