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

func (c *InMemoryVolConfig) Get(name string) (*Volume, error) {
	v := c.Data[name]
	if v == nil {
		return nil, errors.New("no volume found")
	}
	return v, nil
}

func (c *InMemoryVolConfig) Put(name string, volume *Volume) error {
	c.Data[name] = volume
	return nil
}

func (c *InMemoryVolConfig) Remove(name string) error {
	delete(c.Data, name)
	return nil
}

func (c *InMemoryVolConfig) List() (vs []*Volume) {
	for _, v := range c.Data {
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

func (c *FileVolConfig) Get(name string) (*Volume, error) {
	return c.inMemoryVolConfig.Get(name)
}

func (c *FileVolConfig) Put(name string, volume *Volume) error {
	return c.inMemoryVolConfig.Put(name, volume)
}

func (c *FileVolConfig) Remove(name string) error {
	return c.inMemoryVolConfig.Remove(name)
}

func (c *FileVolConfig) List() (vs []*Volume) {
	return c.inMemoryVolConfig.List()
}

func (c *FileVolConfig) Save() error {
	out, err := json.Marshal(c.inMemoryVolConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.filename, out, 0664)
	if err != nil {
		return err
	}
	return nil
}

func (c *FileVolConfig) Load() error {
	m := InMemoryVolConfig{}
	in, err := ioutil.ReadFile(c.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(in, &m)
	if err != nil {
		return err
	}
	c.inMemoryVolConfig = &m
	return nil
}
