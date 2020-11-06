package easycache

import (
	"golang.org/x/sync/singleflight"
	"sync"
)

/**
Easy Cache
Author: https://github.com/amhr
*/
type EasyCache struct {
	layers        []CacheLayer
	resources     map[string]Resource
	resourceLock  sync.RWMutex
	keyGenerator  func(slug string, params ...string) string
	resourceGroup *singleflight.Group
}

// easy cache constructor
func NewEasyCache() *EasyCache {
	var requestGroup singleflight.Group
	return &EasyCache{
		keyGenerator:  defaultKeyGenerator,
		resources:     map[string]Resource{},
		resourceGroup: &requestGroup,
	}
}

// key generator is a function that generates key
// default key generator is located under ./key.go
func (p *EasyCache) SetKeyGenerator(kg func(slug string, params ...string) string) {
	p.keyGenerator = kg
}

// add new cache layers
// cache layers must implement CacheLayer interface
func (p *EasyCache) AddLayer(layers ...CacheLayer) {
	p.layers = append(p.layers, layers...)
}

func (p *EasyCache) GetLayer(layer int) (CacheLayer, bool) {
	if layer >= len(p.layers) {
		return nil, false
	}
	return p.layers[layer], true
}

// handling data fetching
// slug is the key deference between resources
// ex: getUser / getUserByID / getUserByEmail
func (p *EasyCache) AddResource(slug string, resource Resource) {
	p.resourceLock.Lock()
	defer p.resourceLock.Unlock()

	p.resources[slug] = resource
}

func (p *EasyCache) GetResource(slug string) (Resource, bool) {
	p.resourceLock.RLock()
	defer p.resourceLock.RUnlock()
	val, e := p.resources[slug]
	return val, e
}

// provide data
func (p *EasyCache) Provide(slug string, params ...string) ([]byte, error) {
	r, e := p.GetResource(slug)
	if !e {
		return nil, ResourceNotFound{slug: slug}
	}

	// first try in cache
	b, err := p.provideByCache(r, slug, params...)
	// data has been founded in cache
	if err == nil {
		return b, err
	}
	// we didn't find data
	// handling error if it's fatal or we can skip it
	switch err.(type) {
	case ResourceLayerUndefined:
		return b, err
	}
	// check not founded!
	// we couldn't find data in cache
	// using singleflight to handle thundering herds
	v, err, _ := p.resourceGroup.Do(p.keyGenerator(slug, params...), func() (interface{}, error) {
		return p.provideResource(r, slug, params...)
	})
	b = v.([]byte)
	// we had error while providing!
	if err != nil {
		return b, err
	}
	// data is provided, lets set data in cache
	p.Set(b, slug, params...)
	return b, err

}

// provide data from resource
func (p *EasyCache) provideResource(r Resource, slug string, params ...string) ([]byte, error) {
	b, err := r.Provider(slug, params...)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// provide if data is in cache
func (p *EasyCache) provideByCache(r Resource, slug string, params ...string) ([]byte, error) {
	// for each on layers
	for layerIndex := range r.Layers() {
		layer, e := p.GetLayer(layerIndex)
		if !e {
			return nil, ResourceLayerUndefined{
				slug:  slug,
				layer: layerIndex,
			}
		}
		// if this layer has the data returns it
		b, err := layer.Get(p.keyGenerator(slug, params...))
		if err == nil {
			return b, nil
		}
		// otherwise checks other layers
	}
	return nil, &ResourceDidntFindInCache{slug: slug}
}

// set value
func (p *EasyCache) Set(value []byte, slug string, params ...string) error {
	r, e := p.GetResource(slug)
	if !e {
		return ResourceNotFound{slug: slug}
	}
	for layerIndex := range r.Layers() {
		layer, e := p.GetLayer(layerIndex)
		if !e {
			continue
		}
		layer.Set(p.keyGenerator(slug, params...), value)
	}
	return nil
}
