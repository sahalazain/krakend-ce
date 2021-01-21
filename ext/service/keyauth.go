package service

import (
	"errors"
	"time"

	cache "github.com/devopsfaith/krakend-ce/ext/cache"
)

//KeyAuthResponse OPA response model
type KeyAuthResponse struct {
	Result string `json:"result,omitempty" mapstructure:"result"`
}

//KeyAuth keyAuth service interface
type KeyAuth interface {
	Validate(key Cacheable) (string, error)
}

//HTTPKeyAuth http keyAuth service
type HTTPKeyAuth struct {
	address      string
	basePath     string
	cache        cache.Local
	responsePath string
}

//DummyKeyAuth dummy key auth service
type DummyKeyAuth struct {
	Result string
	Error  error
}

//NewHTTPKeyAuth create instance of http keyAuth service
func NewHTTPKeyAuth(address, basePath, responsePath string, cacheDuration, cacheSize int) *HTTPKeyAuth {
	var c cache.Local
	if cacheSize > 0 {
		c, _ = cache.NewLRU(cacheSize)
	}

	if c == nil {
		c = cache.NewMemoryCache(time.Duration(cacheDuration) * time.Second)
	}

	return &HTTPKeyAuth{
		address:      address,
		basePath:     basePath,
		cache:        c,
		responsePath: responsePath,
	}
}

//NewDummyKeyAuth create new dummy auth instance
func NewDummyKeyAuth() *DummyKeyAuth {
	return &DummyKeyAuth{}
}

//Validate validate key api
func (h *HTTPKeyAuth) Validate(key Cacheable) (string, error) {
	hs := key.Hash()

	if rsp, ok := h.cache.Get(hs); ok {
		return rsp.(string), nil
	}

	var rsp map[string]interface{}
	if err := post(h.address, h.basePath, key, &rsp); err != nil {
		return "", err
	}

	res, ok := lookup(h.responsePath, rsp)
	if !ok {
		return "", errors.New("No value within result path")
	}

	id, ok := res.(string)
	if !ok {
		return "", errors.New("Invalid result type")
	}

	h.cache.Set(hs, id)

	return id, nil
}

//Validate validate key api
func (d *DummyKeyAuth) Validate(key Cacheable) (string, error) {
	return d.Result, d.Error
}
