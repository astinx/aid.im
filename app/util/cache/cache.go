package cache

import (
	"fmt"
	"time"
)

var c Cache

type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// start gc routine based on config string settings.
	StartAndGC(config string) error
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func InitCache(adapterName, config string) (err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	c = instanceFunc()
	if err = c.StartAndGC(config); err != nil {
		c = nil
	}
	return
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName, config string) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	if err = adapter.StartAndGC(config); err != nil {
		adapter = nil
	}
	return
}

func Get(key string) interface{} {
	return c.Get(key)
}

func GetMulti(keys []string) []interface{} {
	return c.GetMulti(keys)
}
func Put(key string, val interface{}, timeout time.Duration) error {
	return c.Put(key, val, timeout)
}

// delete cached value by key.
func Delete(key string) error {
	return c.Delete(key)
}

// increase cached int value by key, as a counter.
func Incr(key string) error {
	return c.Incr(key)
}

// decrease cached int value by key, as a counter.
func Decr(key string) error {
	return c.Decr(key)
}

// check if cached value exists or not.
func IsExist(key string) bool {
	return c.IsExist(key)
}

// clear all cache.
func ClearAll() error {
	return c.ClearAll()
}

// start gc routine based on config string settings.
func StartAndGC(config string) error {
	return c.StartAndGC(config)
}
