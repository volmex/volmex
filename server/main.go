package main

import (
	"github.com/docker/go-plugins-helpers/volume"
	"gitlab.mi.hdm-stuttgart.de/fankhauser/volmex"
)

func main() {
	d := volmex.New()
	h := volume.NewHandler(d)
	err := h.ServeUnix("volmex", 0)
	if err != nil {
		panic(err)
	}
}
