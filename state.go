package volmex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// State holds the drivers volume state / configuration
type State interface {
	Get(name string) (*Volume, error)
	Put(name string, volume *Volume) error
	Remove(name string) error
	List() []*Volume
	Save() error
	Load() error
}

// InMemoryState implements the State interface with in-memory based storage
type InMemoryState struct {
	Data map[string]*Volume
	Mux  sync.Mutex
}

// NewInMemoryState returns a new empty InMemoryState instance
func NewInMemoryState() *InMemoryState {
	return &InMemoryState{
		Data: make(map[string]*Volume, 0),
	}
}

// Get returns either the volume configuration for a given volume name or an error if the volume is unknown
func (s *InMemoryState) Get(name string) (*Volume, error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	v := s.Data[name]
	if v == nil {
		return nil, errors.New("volume not found")
	}
	return v, nil
}

// Put stores a volume configuration for a given volume name
func (s *InMemoryState) Put(name string, volume *Volume) error {
	s.Mux.Lock()
	s.Data[name] = volume
	s.Mux.Unlock()
	return nil
}

// Remove removes a volume configuration for a given volume name
func (s *InMemoryState) Remove(name string) error {
	s.Mux.Lock()
	delete(s.Data, name)
	s.Mux.Unlock()
	return nil
}

// List lists all known volume configurations
func (s *InMemoryState) List() (vs []*Volume) {
	s.Mux.Lock()
	for _, v := range s.Data {
		vs = append(vs, v)
	}
	s.Mux.Unlock()
	return vs
}

// Save is an empty implementation (for the State interface)
func (s *InMemoryState) Save() error {
	return nil
}

// Load is an empty implementation (for the State interface)
func (s *InMemoryState) Load() error {
	return nil
}

// FileState implements the State interface with file based storage (by reusing InMemoryState)
type FileState struct {
	inMemoryState *InMemoryState
	filename      string
}

// NewFileState returns a new empty NewFileState instance
func NewFileState(filename string) *FileState {
	return &FileState{
		inMemoryState: NewInMemoryState(),
		filename:      filename,
	}
}

// Get proxies to InMemoryState's Get
func (s *FileState) Get(name string) (*Volume, error) {
	return s.inMemoryState.Get(name)
}

// Put proxies to InMemoryState's Put
func (s *FileState) Put(name string, volume *Volume) error {
	return s.inMemoryState.Put(name, volume)
}

// Remove proxies to InMemoryState's Remove
func (s *FileState) Remove(name string) error {
	return s.inMemoryState.Remove(name)
}

// List proxies to InMemoryState's List
func (s *FileState) List() (vs []*Volume) {
	return s.inMemoryState.List()
}

// Save dumps FileState.inMemoryState to FileState.filename
func (s *FileState) Save() error {
	s.inMemoryState.Mux.Lock()
	defer s.inMemoryState.Mux.Unlock()
	out, err := json.Marshal(s.inMemoryState.Data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.filename, out, 0664)
	if err != nil {
		return err
	}
	return nil
}

// Load fills FileState.inMemoryState with FileState.filename's parsed content (if present)
func (s *FileState) Load() error {
	s.inMemoryState.Mux.Lock()
	defer s.inMemoryState.Mux.Unlock()
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return nil
	}
	d := make(map[string]*Volume, 0)
	in, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return fmt.Errorf("could not parse state from %v: %v", s.filename, err)
	}
	err = json.Unmarshal(in, &d)
	if err != nil {
		return err
	}
	s.inMemoryState.Data = d
	return nil
}
