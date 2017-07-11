package volmex

import (
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-plugins-helpers/volume"
)

func TestFileState(t *testing.T) {
	fname := fmt.Sprintf("/tmp/test.json %v", time.Now().Unix())

	// create a FileState, add a volume configuration and save the state
	s := NewFileState(fname)
	s.Put("foo", &Volume{
		Volume: volume.Volume{
			Name: "foo",
		},
	})
	err := s.Save()
	if err != nil {
		t.Fatal(err)
	}

	// create new FileState, load the state and check if our previously added configuration is still there
	ss := NewFileState(fname)
	err = ss.Load()
	if err != nil {
		t.Fatal(err)
	}
	v, err := ss.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if v.Volume.Name != "foo" {
		t.Fatal("did not load volume configuration correctly")
	}
}
