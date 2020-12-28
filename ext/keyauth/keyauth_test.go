package keyauth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/devopsfaith/krakend-ce/ext/service"
	"github.com/stretchr/testify/assert"
)

func TestBodyRequest(t *testing.T) {
	const json = `{"key_api": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","id":10001}`
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "body.key_api",
		IDHeaderName:   defaultHeaderName,
	}

	ds := service.NewDummyKeyAuth()

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(json)))

	k, err := cfg.extractKey(req)
	assert.Nil(t, err)
	assert.Equal(t, k, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")

	cfg.Service = ds
	ds.Error = nil
	ds.Result = "partner1"

	r, err := cfg.validateKey(req)
	assert.Nil(t, err)
	assert.True(t, r)
	fmt.Println(cfg.IDHeaderName)
	assert.Equal(t, req.Header.Get(cfg.IDHeaderName), "partner1")

	ds.Result = ""
	r, err = cfg.validateKey(req)
	assert.False(t, r)
}

func TestHeaderRequest(t *testing.T) {
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "header.X-API-Key",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	req.Header.Add("X-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")

	k, err := cfg.extractKey(req)
	assert.Nil(t, err)
	assert.Equal(t, k, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")

}

func TestQueryRequest(t *testing.T) {
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "query.api-key",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?api-key=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")

	k, err := cfg.extractKey(req)
	assert.Nil(t, err)
	assert.Equal(t, k, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")

}

func TestInvalidRequest(t *testing.T) {
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "jwt.apikey",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")

	k, err := cfg.extractKey(req)
	assert.NotNil(t, err)
	assert.Equal(t, k, "")

}

func TestInvalidBodyRequest(t *testing.T) {
	const json = `{"key_api": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9","id":10001}`
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "body.keyapi",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(json)))

	k, err := cfg.extractKey(req)
	assert.NotNil(t, err)
	assert.Equal(t, k, "")
}

func TestInvalidHeaderRequest(t *testing.T) {
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "header.X-API-Key",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	req.Header.Add("XAPI-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")

	k, err := cfg.extractKey(req)
	assert.NotNil(t, err)
	assert.Equal(t, k, "")

}

func TestInvalidQueryRequest(t *testing.T) {
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		KeyPath:        "query.apikey",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?api-key=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")

	k, err := cfg.extractKey(req)
	assert.NotNil(t, err)
	assert.Equal(t, k, "")

}
