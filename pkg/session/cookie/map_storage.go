package cookie

import "github.com/gomodule/redigo/redis"

type mapStorage struct {
	storage map[string][]byte
}

func NewMapStorage(conn redis.Conn, ttl int) *mapStorage {
	return &mapStorage{storage: make(map[string][]byte)}
}

func (s *mapStorage) Add(mkey string, serialized []byte) error {
	s.storage[mkey] = serialized
	return nil
}

func (s *mapStorage) Get(mkey string) ([]byte, error) {
	return s.storage[mkey], nil
}

func (s *mapStorage) GetTTL() int {
	return 0
}
