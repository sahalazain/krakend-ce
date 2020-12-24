package service

import (
	"strings"
	"time"

	cache "github.com/devopsfaith/krakend-ce/ext/cache"
)

//Cacheable cacheable object
type Cacheable interface {
	Hash() [32]byte
}

//Policy policy service interface
type Policy interface {
	Evaluate(pkg, directive string, data Cacheable) (bool, error)
}

//HTTPOPA http OPA service
type HTTPOPA struct {
	address  string
	basePath string
	cache    *cache.MemoryCache
}

//DummyOPA dummy OPA service
type DummyOPA struct {
	Result bool
	Error  error
}

//Response OPA response model
type Response struct {
	Result bool `json:"result,omitempty" mapstructure:"result"`
}

//NewHTTPOPA create new http OPA service instance
func NewHTTPOPA(address, basePath string, cacheDuration int) *HTTPOPA {
	return &HTTPOPA{
		address:  address,
		basePath: basePath,
		cache:    cache.NewMemoryCache(time.Duration(cacheDuration) * time.Second),
	}
}

//NewDummyOPA create new dummy OPA service instance
func NewDummyOPA() *DummyOPA {
	return &DummyOPA{}
}

//Evaluate evaluate input request against policy
func (d *DummyOPA) Evaluate(pkg, directive string, data Cacheable) (bool, error) {
	return d.Result, d.Error
}

//Evaluate evaluate input request against policy
func (h *HTTPOPA) Evaluate(pkg, directive string, data Cacheable) (bool, error) {

	hs := data.Hash()

	if rsp, ok := h.cache.Get(hs); ok {
		return rsp.(bool), nil
	}

	path := h.basePath + strings.ReplaceAll(pkg, ".", "/") + "/" + directive
	var rsp Response
	if err := post(h.address, path, data, &rsp); err != nil {
		return false, err
	}

	h.cache.Set(hs, rsp.Result)

	return rsp.Result, nil
}
