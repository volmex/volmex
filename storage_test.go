package volmex

import (
	"testing"
	"github.com/docker/go-plugins-helpers/volume"
)

func TestFileStorage(t *testing.T) {
	s := NewFileStorage("/tmp/test.json")

	s.Put("foo", &Volume{
		Volume: volume.Volume{
			Name: "foo",
	}})

	err := s.Save()

	if err != nil {
		t.Fatal(err)
	}

	ss := NewFileStorage("/tmp/test.json")
	err = ss.Load()
	if err != nil {
		t.Fatal(err)
	}

	v, err := ss.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if v.Volume.Name != "foo" {
		t.Fatal("did not load volumes correctly")
	}
}