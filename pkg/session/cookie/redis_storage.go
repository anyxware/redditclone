package cookie

import "github.com/gomodule/redigo/redis"

type redisStorage struct {
	conn redis.Conn
	ttl  int
}

func NewRedisStorage(conn redis.Conn, ttl int) *redisStorage {
	return &redisStorage{conn: conn, ttl: ttl}
}

func (s *redisStorage) Add(mkey string, serialized []byte) error {
	_, err := s.conn.Do("SET", mkey, serialized, "EX", s.ttl)
	if err != nil {
		return err
	}
	return nil
}

func (s *redisStorage) Get(mkey string) ([]byte, error) {
	data, err := s.conn.Do("GET", mkey)
	if err != nil {
		return nil, err
	}
	return data.([]byte), nil
}

func (s *redisStorage) GetTTL() int {
	return s.ttl
}
