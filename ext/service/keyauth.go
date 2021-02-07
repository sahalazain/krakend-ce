package service

import (
	"time"

	cache "github.com/devopsfaith/krakend-ce/ext/cache"
)

//KeyAuthResponse OPA response model
type KeyAuthResponse struct {
	Result string `json:"result,omitempty" mapstructure:"result"`
}

//KeyAuth keyAuth service interface
type KeyAuth interface {
	Validate(key Cacheable) (map[string]interface{}, error)
}

//HTTPKeyAuth http keyAuth service
type HTTPKeyAuth struct {
	address  string
	basePath string
	cache    cache.Local
}

//DummyKeyAuth dummy key auth service
type DummyKeyAuth struct {
	Result map[string]interface{}
	Error  error
}

//NewHTTPKeyAuth create instance of http keyAuth service
func NewHTTPKeyAuth(address, basePath string, cacheDuration, cacheSize int) *HTTPKeyAuth {
	var c cache.Local
	if cacheSize > 0 {
		c, _ = cache.NewLRU(cacheSize)
	}

	if c == nil {
		c = cache.NewMemoryCache(time.Duration(cacheDuration) * time.Second)
	}

	return &HTTPKeyAuth{
		address:  address,
		basePath: basePath,
		cache:    c,
	}
}

//NewDummyKeyAuth create new dummy auth instance
func NewDummyKeyAuth() *DummyKeyAuth {
	return &DummyKeyAuth{}
}

//Validate validate key api
func (h *HTTPKeyAuth) Validate(key Cacheable) (map[string]interface{}, error) {
	hs := key.Hash()

	if rsp, ok := h.cache.Get(hs); ok {
		return rsp.(map[string]interface{}), nil
	}

	var rsp map[string]interface{}
	if err := post(h.address, h.basePath, key, &rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

//Validate validate key api
func (d *DummyKeyAuth) Validate(key Cacheable) (map[string]interface{}, error) {
	return d.Result, d.Error
}
