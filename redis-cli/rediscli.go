package rediscli

import (
	"context"
	"fmt"
    "github.com/go-redis/redis/v8"
    "bytes"
    "encoding/gob"  
    "log"
    "time"  
)

type Cache struct {
	client *redis.Client

	ctx context.Context
} 

func NewCache() Cache {
	return Cache{ctx:  context.Background()}
}

func (c Cache) cli() *redis.Client {
	if c.client == nil {
		c.client = redis.NewClient(&redis.Options{
	        Addr:     "localhost:6379",
	        Password: "", // no password set
	        DB:       0,  // use default DB
	    })
	}
	
	return c.client
}

func (c Cache) Put(key string, e interface{}) {
    var network bytes.Buffer
    enc := gob.NewEncoder(&network)
    err := enc.Encode(e)
    if err != nil {
        log.Fatal("encode error:", err)
    }
    
    err = c.cli().Set(c.ctx, key, network.Bytes(), 0).Err()
    if err != nil {
        panic(err)
    }
}

// decode value into 2nd argument, returning true when there was one stored for the key
func (c Cache) Get(key string, o interface{}) bool {
	atStart := time.Now()
	value, err := c.cli().Get(c.ctx, key).Bytes()
	if err == redis.Nil {
        return false
    } else if err != nil {
        panic(err)
    }
    
    buf := bytes.NewBuffer(value)
    dec := gob.NewDecoder(buf) 
    // Decode (receive) the value.
    err = dec.Decode(o)
    if err == nil {
    	fmt.Printf("%v for cached key %s\n", (time.Now().Sub(atStart)), key)
    } else {
        log.Fatal("decode error:", err)
    }
    
    return true
}
