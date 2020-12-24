package opa

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/devopsfaith/krakend-ce/ext/service"
	"github.com/stretchr/testify/assert"
)

func TestBasicRequest(t *testing.T) {

	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		PackageName:    "opa.test",
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha/", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")

	oreq := cfg.buildRequest(req)
	assert.NotNil(t, oreq, "Should not nil")
	assert.Equal(t, oreq.Input.Method, "GET", "Method should GET")
	assert.Equal(t, oreq.Input.Path, []string{"echo", "alpha"}, "Path should be [echo, alpha]")
}

func TestWithPayload(t *testing.T) {
	const json = `{"name":{"first":"Janet", "last":"Prichard" },"age":47}`
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		PackageName:    "opa.test",
		PayloadMap: map[string]string{
			"token":      "jwt.raw",
			"auth":       "header.Authorization",
			"username":   "jwt.payload.name",
			"jwt_type":   "jwt.header.typ",
			"age":        "body.age",
			"first_name": "body.name.first",
			"body":       "body.raw",
			"headers":    "header.raw",
			"city":       "query.city",
		},
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")
	req.Header.Add(authHeader, "Bearer "+token)
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(json)))

	oreq := cfg.buildRequest(req)
	assert.NotNil(t, oreq, "Should not nil")
	assert.Equal(t, oreq.Input.Payload["token"], token, "Should be raw token")
	assert.Equal(t, oreq.Input.Payload["auth"], "Bearer "+token, "Should be raw auth header")
	assert.Equal(t, oreq.Input.Payload["username"], "John Doe", "Should be John Doe")
	assert.Equal(t, oreq.Input.Payload["jwt_type"], "JWT", "Should be JWT")
	assert.Equal(t, oreq.Input.Payload["age"], float64(47), "Should be 47")
	assert.Equal(t, oreq.Input.Payload["first_name"], "Janet", "Should be Janet")
	assert.Equal(t, oreq.Input.Payload["city"], "Jakarta", "Should be Jakarta")
}

func TestCheckPermission(t *testing.T) {
	const json = `{"name":{"first":"Janet", "last":"Prichard" },"age":47}`
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		PackageName:    "opa.test",
		PayloadMap: map[string]string{
			"token":      "jwt.raw",
			"auth":       "header.Authorization",
			"username":   "jwt.payload.name",
			"jwt_type":   "jwt.header.typ",
			"age":        "body.age",
			"first_name": "body.name.first",
			"body":       "body.raw",
			"headers":    "header.raw",
			"city":       "query.city",
		},
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")
	req.Header.Add(authHeader, "Bearer "+token)
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(json)))

	ds := service.NewDummyOPA()
	cfg.Service = ds
	ds.Error = nil
	ds.Result = true

	rsp, err := cfg.checkPermission(req)
	assert.Nil(t, err)
	assert.True(t, rsp)

	ds.Result = false
	rs, err := cfg.checkPermission(req)
	assert.False(t, rs)

}
