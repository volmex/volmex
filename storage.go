package volmex

import "errors"

type Storage interface {
	Get(name string) (*Volume, error)
	Put(name string, volume *Volume) error
	Remove(name string) error
	List() []*Volume
}

type InMemoryStorage struct {
	data map[string]*Volume
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]*Volume, 0),
	}
}

func (s *InMemoryStorage) Get(name string) (*Volume, error) {
	v := s.data[name]
	if v == nil {
		return nil, errors.New("no volume found")
	}
	return v, nil
}

func (s *InMemoryStorage) Put(name string, volume *Volume) error {
	s.data[name] = volume
	return nil
}

func (s *InMemoryStorage) Remove(name string) error {
	delete(s.data, name)
	return nil
}

func (s *InMemoryStorage) List() (vs []*Volume) {
	for _, v := range s.data {
		vs = append(vs, v)
	}
	return vs
}
