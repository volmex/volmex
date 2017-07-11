package volmex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type State interface {
	Get(name string) (*Volume, error)
	Put(name string, volume *Volume) error
	Remove(name string) error
	List() []*Volume
	Save() error
	Load() error
}

type InMemoryState struct {
	Data map[string]*Volume
	Mux  sync.Mutex
}

func NewInMemoryState() *InMemoryState {
	return &InMemoryState{
		Data: make(map[string]*Volume, 0),
	}
}

func (s *InMemoryState) Get(name string) (*Volume, error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	v := s.Data[name]
	if v == nil {
		return nil, errors.New("no volume found")
	}
	return v, nil
}

func (s *InMemoryState) Put(name string, volume *Volume) error {
	s.Mux.Lock()
	s.Data[name] = volume
	s.Mux.Unlock()
	return nil
}

func (s *InMemoryState) Remove(name string) error {
	s.Mux.Lock()
	delete(s.Data, name)
	s.Mux.Unlock()
	return nil
}

func (s *InMemoryState) List() (vs []*Volume) {
	s.Mux.Lock()
	for _, v := range s.Data {
		vs = append(vs, v)
	}
	s.Mux.Unlock()
	return vs
}

func (s *InMemoryState) Save() error {
	return nil
}

func (s *InMemoryState) Load() error {
	return nil
}

type FileState struct {
	inMemoryState *InMemoryState
	filename      string
}

func NewFileState(filename string) *FileState {
	return &FileState{
		inMemoryState: NewInMemoryState(),
		filename:      filename,
	}
}

func (s *FileState) Get(name string) (*Volume, error) {
	return s.inMemoryState.Get(name)
}

func (s *FileState) Put(name string, volume *Volume) error {
	return s.inMemoryState.Put(name, volume)
}

func (s *FileState) Remove(name string) error {
	return s.inMemoryState.Remove(name)
}

func (s *FileState) List() (vs []*Volume) {
	return s.inMemoryState.List()
}

func (s *FileState) Save() error {
	s.inMemoryState.Mux.Lock()
	defer s.inMemoryState.Mux.Unlock()
	out, err := json.Marshal(s.inMemoryState)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.filename, out, 0664)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileState) Load() error {
	s.inMemoryState.Mux.Lock()
	defer s.inMemoryState.Mux.Unlock()
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return nil
	}
	m := InMemoryState{}
	in, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return fmt.Errorf("could not load state from %v: %v", s.filename, err)
	}
	err = json.Unmarshal(in, &m)
	if err != nil {
		return err
	}
	s.inMemoryState = &m
	return nil
}
