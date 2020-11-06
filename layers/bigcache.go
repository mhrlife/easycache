package layers

import bc "github.com/allegro/bigcache"

type BigCache struct {
	Cache *bc.BigCache
}

func (b *BigCache) Get(key string) ([]byte, error) {
	return b.Cache.Get(key)
}

func (b *BigCache) Set(key string, value []byte) error {
	return b.Cache.Set(key, value)
}
