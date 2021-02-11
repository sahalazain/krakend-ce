package jwtmap

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/devopsfaith/krakend/config"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestJWTMap(t *testing.T) {
	cfg := configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"jwt_map": map[string]interface{}{
				"header.name":     "payload.name",
				"query.sub":       "payload.sub",
				"query.name":      "payload.name",
				"body.name.first": "payload.name",
			},
		},
	})

	assert.NotNil(t, cfg)

	const jsonData = `{"name":{"first":"Janet", "last":"Prichard" },"age":47}`
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	req, err := http.NewRequest("POST", "http://localhost:8000/echo?name=John", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req, "Should not nil")
	req.Header.Add(authHeader, "Bearer "+token)
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(jsonData)))

	extractClaim(cfg, req)

	assert.Equal(t, "John Doe", req.Header.Get("name"))
	assert.Equal(t, "1234567890", req.URL.Query().Get("sub"))
	assert.Equal(t, "John Doe", req.URL.Query().Get("name"))

	raw, _ := ioutil.ReadAll(req.Body)
	var body map[string]interface{}
	json.Unmarshal(raw, &body)

	val := gjson.Get(string(raw), "name.first")
	assert.Equal(t, "John Doe", val.Str)
}
