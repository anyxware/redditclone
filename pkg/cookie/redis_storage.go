package cookie

import "github.com/gomodule/redigo/redis"

type redisStorage struct {
	conn redis.Conn
}

func NewRedisStorage(conn redis.Conn) *redisStorage {
	return &redisStorage{conn: conn}
}

func (s *redisStorage) Add(mkey string, serialized []byte) error {
	_, err := s.conn.Do("SET", mkey, serialized, "EX", 86400)
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
