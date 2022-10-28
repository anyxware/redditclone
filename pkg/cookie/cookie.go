package cookie

type storage interface {
	Add(mkey string, serialized []byte) error
	Get(mkey string) ([]byte, error)
}

type Manager struct {
	storage storage
}

func NewManager(storage storage) Manager {
	return Manager{storage: storage}
}

func (m Manager) AddCookie(mkey string, serialized []byte) error {
	return m.storage.Add(mkey, serialized)
}

func (m Manager) GetCookie(mkey string) ([]byte, error) {
	return m.storage.Get(mkey)
}