package volmex

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type VolConfig interface {
	Get(name string) (*Volume, error)
	Put(name string, volume *Volume) error
	Remove(name string) error
	List() []*Volume
}

type InMemoryVolConfig struct {
	Data map[string]*Volume
}

func NewInMemoryVolConfig() *InMemoryVolConfig {
	return &InMemoryVolConfig{
		Data: make(map[string]*Volume, 0),
	}
}

func (s *InMemoryVolConfig) Get(name string) (*Volume, error) {
	v := s.Data[name]
	if v == nil {
		return nil, errors.New("no volume found")
	}
	return v, nil
}

func (s *InMemoryVolConfig) Put(name string, volume *Volume) error {
	s.Data[name] = volume
	return nil
}

func (s *InMemoryVolConfig) Remove(name string) error {
	delete(s.Data, name)
	return nil
}

func (s *InMemoryVolConfig) List() (vs []*Volume) {
	for _, v := range s.Data {
		vs = append(vs, v)
	}
	return vs
}

type FileVolConfig struct {
	inMemoryVolConfig *InMemoryVolConfig
	filename          string
}

func NewFileVolConfig(filename string) *FileVolConfig {
	return &FileVolConfig{
		inMemoryVolConfig: NewInMemoryVolConfig(),
		filename:          filename,
	}
}

func (s *FileVolConfig) Get(name string) (*Volume, error) {
	return s.inMemoryVolConfig.Get(name)
}

func (s *FileVolConfig) Put(name string, volume *Volume) error {
	return s.inMemoryVolConfig.Put(name, volume)
}

func (s *FileVolConfig) Remove(name string) error {
	return s.inMemoryVolConfig.Remove(name)
}

func (s *FileVolConfig) List() (vs []*Volume) {
	return s.inMemoryVolConfig.List()
}

func (s *FileVolConfig) Save() error {
	out, err := json.Marshal(s.inMemoryVolConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.filename, out, 0664)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileVolConfig) Load() error {
	m := InMemoryVolConfig{}
	in, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(in, &m)
	if err != nil {
		return err
	}
	s.inMemoryVolConfig = &m
	return nil
}
