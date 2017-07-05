package volmex

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Storage interface {
	Get(name string) (*Volume, error)
	Put(name string, volume *Volume) error
	Remove(name string) error
	List() []*Volume
}

type InMemoryStorage struct {
	Data map[string]*Volume
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Data: make(map[string]*Volume, 0),
	}
}

func (s *InMemoryStorage) Get(name string) (*Volume, error) {
	v := s.Data[name]
	if v == nil {
		return nil, errors.New("no volume found")
	}
	return v, nil
}

func (s *InMemoryStorage) Put(name string, volume *Volume) error {
	s.Data[name] = volume
	return nil
}

func (s *InMemoryStorage) Remove(name string) error {
	delete(s.Data, name)
	return nil
}

func (s *InMemoryStorage) List() (vs []*Volume) {
	for _, v := range s.Data {
		vs = append(vs, v)
	}
	return vs
}

type FileStorage struct {
	inMemoryStorage *InMemoryStorage
	filename        string
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		inMemoryStorage: NewInMemoryStorage(),
		filename:        filename,
	}
}

func (s *FileStorage) Get(name string) (*Volume, error) {
	return s.inMemoryStorage.Get(name)
}

func (s *FileStorage) Put(name string, volume *Volume) error {
	return s.inMemoryStorage.Put(name, volume)
}

func (s *FileStorage) Remove(name string) error {
	return s.inMemoryStorage.Remove(name)
}

func (s *FileStorage) List() (vs []*Volume) {
	return s.inMemoryStorage.List()
}

func (s *FileStorage) Save() error {
	out, err := json.Marshal(s.inMemoryStorage)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.filename, out, 0664)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Load() error {
	m := InMemoryStorage{}
	in, err := ioutil.ReadFile(s.filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(in, &m)
	if err != nil {
		return err
	}
	s.inMemoryStorage = &m
	return nil
}
