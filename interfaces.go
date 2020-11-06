package easycache

type CacheLayer interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

type Resource interface {
	Layers() []int
	Provider(slug string, params ...string) ([]byte, error)
}
