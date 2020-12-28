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
	Validate(key Cacheable) (string, error)
}

//HTTPKeyAuth http keyAuth service
type HTTPKeyAuth struct {
	address  string
	basePath string
	cache    *cache.MemoryCache
}

//DummyKeyAuth dummy key auth service
type DummyKeyAuth struct {
	Result string
	Error  error
}

//NewHTTPKeyAuth create instance of http keyAuth service
func NewHTTPKeyAuth(address, basePath string, cacheDuration int) *HTTPKeyAuth {
	return &HTTPKeyAuth{
		address:  address,
		basePath: basePath,
		cache:    cache.NewMemoryCache(time.Duration(cacheDuration) * time.Second),
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

	var rsp KeyAuthResponse
	if err := post(h.address, h.basePath, key, &rsp); err != nil {
		return "", err
	}

	h.cache.Set(hs, rsp.Result)

	return rsp.Result, nil
}

//Validate validate key api
func (d *DummyKeyAuth) Validate(key Cacheable) (string, error) {
	return d.Result, d.Error
}
