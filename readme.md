# EasyCache

Easy cache provides you a wrapper over Redis and BigCache.
  - You can use available layers (bigcache,redis)
  - You can implement your own caching layers
  - It uses singleflight library which prevents Thundering Herd problem.
### Installation
```
go get -u github.com/mhrlife/easycache
```
### Creating Resource
This library acts like google's GroupCache. You need to implement Resource for each resources.
```
type Resource interface {
	Layers() []int
	Provider(slug string, params ...string) ([]byte, error)
}
```
 - Layers()  Determinates this resource saves data in which layers.
 - Provider()  In this function you need to implement how to provide data.
### Examples
Imagine you want to cache user infos. In this example we cache user infos in redis with the key format userinfo:$id
First we need to create a resource.
// for testing errors
```
type UserInfoResource struct {
    DB *Gorm.DB
}
func (c *UserInfoResource) Provider(slug string, params ...string) ([]byte, error) {
	// provide user from database or somewhere else
	return UserInBytes , nil
}
func (c UserInfoResource) Layers() []int {
	return []int{0} // which is redis, we need to add it when we are going to create a new EasyCache instance.
}
```
After creating resource we need to create an instance of EasyCache.
```
import (
    "github.com/go-redis/redis/v8"
    "github.com/mhrlife/easycache"
    "github.com/mhrlife/easycache/layers"
)
  
    
func main(){
    ec := easycache.NewEasyCache()
    rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
    ec.AddLayer(&layers.Redis{
        Cache: rdb,
        Ttl:   100 * time.Hour, // ttl of cache
    })
    resource := &UserInfoResource{
    	DB: db, // for example gorm.DB
    }
    ec.AddResource("userinfo", resource)
}
```
And now we can easily access to user info with
```
   bytes, err := ec.Provide("userinfo","3") 
```

### Todos

 - Support for other redis commands

License
----

MIT


**Free Software, Hell Yeah!**