package todo

type ToDo struct {
	ID          string
	Title       string
	Description string
}

type KV interface {
	Set(key string, val []byte) error
	Get(key string) ([]byte, bool)
	Del(key string) bool
	Exists(key string) bool
}
