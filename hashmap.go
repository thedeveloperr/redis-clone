package hashmap

import (
	"sync"
	"time"
)

type HashMap interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Expire(key string, timeoutSeconds int) int
}

type Value struct {
	value        string
	setAt        time.Time
	expireAfter  time.Duration
	shouldExpire bool
}

type ConcurrentMap struct {
	mutex sync.RWMutex
	data  map[string]*Value
}

func Create() *ConcurrentMap {
	hashmap := ConcurrentMap{
		data: make(map[string]*Value),
	}
	return &hashmap
}

func (c *ConcurrentMap) Set(key string, value string) {
	c.mutex.Lock()
	c.data[key] = &Value{
		value:        value,
		setAt:        time.Now(),
		expireAfter:  0,
		shouldExpire: false,
	}
	c.mutex.Unlock()
}

func (c *ConcurrentMap) Get(key string) (string, bool) {
	c.mutex.RLock()
	valueItem, exists := c.data[key]
	if !exists {
		c.mutex.RUnlock()
		return "", false
	}

	// To improve accuracy of EXPIRE, in case time.AfterFunc runs later and a get call is made earlier
	if valueItem.shouldExpire && time.Now().Sub(valueItem.setAt) > valueItem.expireAfter {
		c.mutex.RUnlock()
		return "", false
	}
	c.mutex.RUnlock()
	return valueItem.value, exists
}

func (c *ConcurrentMap) Expire(key string, timeoutSeconds int) int {
	if val, ok := c.Get(key); !ok {
		_ = val
		return 0
	}
	c.mutex.Lock()
	c.data[key].shouldExpire = true
	c.data[key].expireAfter = time.Duration(timeoutSeconds) * time.Second
	c.mutex.Unlock()
	time.AfterFunc(time.Duration(timeoutSeconds)*time.Second, func() {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
	})
	return 1
}
