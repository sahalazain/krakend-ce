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

func TestJWTPayload(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjN1TVJqY1BtZEJLTE9nRlczYVJzRU1YQ0VXaWllVWduT0FYVDBJWG04Rm89In0.eyJhdWQiOiJhcHBzIiwiZXhwIjoxNjExMDQzMDkzLCJncm91cHMiOlsiZGVmYXVsdCJdLCJpc3MiOiJzaWNlcGF0IiwicG9saWN5IjoiZGVmYXVsdCIsInJvbGVzIjpbImRlZmF1bHQiXSwic3ViIjoiR1g2dnZkZzEyRHdubldNbkRkYzV0TWV1MlFHeTljOUxOZlBpMTl0N2J3Vm0ifQ.fUqK8JfGyl8etfpsEZZ9K49D89iycINZ-gq_E1stY87aoTeRFZzTaFwGASyslL8sRm6WRvGE79Tg8CkU3blTt3_Ngl5j41pZ__aL-uK92kgL6aSsquFSZG8XLJxpY8TZPN1PNuOiWGQ3LDTqBozeiqEGaX38LMNSQ4CAGLZS9-wCN3oUdRR-V_ulfdAOTcBSbImtJY6JdamXv_3UhsvjVl5ExQOFIKicY9KyaCLDx9oghQLNvgJcLaFm9z8jYHPIjA1IrAe84-2eWNSwFU9XcGwkpoOr1G4QfX2Cs2FoKWCR_sqQg4qMZRdlDgZFYwkv4ZmXN-hdUmS7OH73LckwmA"
	cfg := &xtraConfig{
		ServiceAddress: "http://localhost:8080",
		PackageName:    "opa.test",
		PayloadMap: map[string]string{
			"subject": "jwt.payload.sub",
			"groups":  "jwt.payload.groups",
		},
	}

	assert.NotNil(t, cfg)

	req, err := http.NewRequest("GET", "http://localhost:8000/echo/alpha?city=Jakarta", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")
	req.Header.Add(authHeader, "Bearer "+token)

	oreq := cfg.buildRequest(req)
	assert.NotNil(t, oreq, "Should not nil")
	assert.Equal(t, "GX6vvdg12DwnnWMnDdc5tMeu2QGy9c9LNfPi19t7bwVm", oreq.Input.Payload["subject"])
	assert.Equal(t, []string{"default"}, oreq.Input.Payload["groups"])
}
