package easycache

import (
	"github.com/allegro/bigcache"
	"github.com/go-redis/redis/v8"
	"github.com/mhrlife/easycache/layers"
	"strings"
	"testing"
	"time"
)

func GetCacheWithBigCache(t *testing.T) *EasyCache {
	ec := NewEasyCache()
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		t.Fatalf("big cache external construction %v", err)
	}
	ec.AddLayer(&layers.BigCache{Cache: cache})
	return ec
}

func GetCacheWithRedis(t *testing.T) *EasyCache {
	ec := NewEasyCache()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ec.AddLayer(&layers.Redis{
		Cache: rdb,
		Ttl:   100 * time.Second,
	})
	return ec
}

func TestBigCacheLayer(t *testing.T) {
	ec := GetCacheWithBigCache(t)
	bg, _ := ec.GetLayer(0)
	err := bg.Set("Hello", []byte("Hi"))
	if err != nil {
		t.Fatalf("big cache set error %v", err)
	}
	b, err := bg.Get("Hello")
	if err != nil {
		t.Fatalf("big cache get error %v", err)
	}
	if string(b) != "Hi" {
		t.Fatalf("big cache get/set not working")
	}
}

func TestRedis(t *testing.T) {
	ec := GetCacheWithRedis(t)
	bg, _ := ec.GetLayer(0)
	err := bg.Set("Hello", []byte("Hi"))
	if err != nil {
		t.Fatalf("redis set error %v", err)
	}
	b, err := bg.Get("Hello")
	if err != nil {
		t.Fatalf("redis get error %v", err)
	}
	if string(b) != "Hi" {
		t.Fatalf("redis get/set not working")
	}
}

// test resource and provide
func TestResource(t *testing.T) {
	ec := GetCacheWithBigCache(t)
	resource := &customResource{
		counter: 0,
	}
	ec.AddResource("getUser", resource)
	b, err := ec.Provide("getUser", "2", "3")
	if err != nil {
		t.Fatalf("error while providing %v", err)
	}
	if string(b) != "getUser:2-3" {
		t.Fatalf("%s != %s", string(b), "getUser:2-3")
	}

	ec.Provide("getUser", "2", "3")
	ec.Provide("getUser", "2", "3")
	ec.Provide("getUser", "2", "3")
	// test for provide calls
	if resource.counter != 1 {
		t.Fatalf("too many provide calls. (%d)", resource.counter)
	}
}

// test resource and provide
func TestUpdate(t *testing.T) {
	ec := GetCacheWithBigCache(t)
	resource := &customResource{
		counter: 0,
	}
	ec.AddResource("getUser", resource)
	b, err := ec.Provide("getUser", "2", "3")
	if err != nil {
		t.Fatalf("error while providing %v", err)
	}
	if string(b) != "getUser:2-3" {
		t.Fatalf("%s != %s", string(b), "getUser:2-3")
	}

	err = ec.Set([]byte("updated"), "getUser", "2", "3")
	if err != nil {
		t.Fatalf("error while setting value %v", err)
	}
	b, err = ec.Provide("getUser", "2", "3")
	if err != nil {
		t.Fatalf("error while providing %v", err)
	}

	if string(b) != "updated" {
		t.Fatalf("%s != %s", string(b), "getUser:2-3")
	}

}

// Test resource cannot provide
func TestResourceError(t *testing.T) {
	ec := GetCacheWithBigCache(t)
	resource := &resourceNotFound{}
	ec.AddResource("getUser", resource)
	_, err := ec.Provide("getUser", "2", "3")
	if _, ok := err.(CannotProvide); ok {

	} else {
		t.Fatalf("type is not cannot provide")
	}
}

// test resource and provide
func TestResourceWithThunderherd(t *testing.T) {
	ec := GetCacheWithBigCache(t)
	resource := &resourceWithDelay{
		counter: 0,
	}
	ec.AddResource("getUser", resource)
	go ec.Provide("getUser", "2", "3")
	go ec.Provide("getUser", "2", "3")
	go ec.Provide("getUser", "2", "3")
	go ec.Provide("getUser", "3", "3")
	go ec.Provide("getUser", "3", "3")
	go ec.Provide("getUser", "3", "3")
	time.Sleep(150 * time.Millisecond)
	// test for provide calls
	if resource.counter != 2 {
		t.Fatalf("too many provide calls. (%d)", resource.counter)
	}

}

type customResource struct {
	counter int
}

func (c *customResource) Provider(slug string, params ...string) ([]byte, error) {
	c.counter++
	return []byte(slug + ":" + strings.Join(params, "-")), nil
}

func (c customResource) Layers() []int {
	return []int{0}
}

// for testing thunder herd
type resourceWithDelay struct {
	counter int
}

func (c *resourceWithDelay) Provider(slug string, params ...string) ([]byte, error) {
	time.Sleep(100 * time.Millisecond)
	c.counter++
	return []byte(slug + ":" + strings.Join(params, "-")), nil
}

func (c resourceWithDelay) Layers() []int {
	return []int{0}
}

// for testing errors
type resourceNotFound struct {
}

func (c *resourceNotFound) Provider(slug string, params ...string) ([]byte, error) {
	return nil, CannotProvide{}
}

func (c resourceNotFound) Layers() []int {
	return []int{0}
}
