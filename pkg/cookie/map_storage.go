package cookie

type mapStorage struct {
	storage map[string][]byte
}

func NewMapStorage() *mapStorage {
	return &mapStorage{storage: make(map[string][]byte)}
}

func (s *mapStorage) Add(mkey string, serialized []byte) error {
	s.storage[mkey] = serialized
	return nil
}

func (s *mapStorage) Get(mkey string) ([]byte, error) {
	return s.storage[mkey], nil
}
