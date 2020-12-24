package opa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"net/http"
	"strings"

	"crypto/sha256"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//Request OPA request model
type Request struct {
	Input Input `json:"input,omitempty" mapstructure:"input"`
}

//Hash calculate request hash
func (r *Request) Hash() [32]byte {
	val := fmt.Sprintf("%v", r)
	return sha256.Sum256([]byte(val))
}

//Response OPA response model
type Response struct {
	Result bool `json:"result,omitempty" mapstructure:"result"`
}

//Input OPA input model
type Input struct {
	Method  string                 `json:"method,omitempty" mapstructure:"method"`
	Path    []string               `json:"path,omitempty" mapstructure:"path"`
	Payload map[string]interface{} `json:"payload,omitempty" mapstructure:"payload"`
}

//HandlerFactory Open Policy Agent handler factory
func HandlerFactory(l logging.Logger, next krakendgin.HandlerFactory) krakendgin.HandlerFactory {
	//l.Debug("Enabling OPA handler ")
	return func(remote *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		handlerFunc := next(remote, p)

		conf := configGetter(remote.ExtraConfig)

		if conf == nil {
			//l.Debug("[OPA] No config for policy agent ")
			return func(c *gin.Context) {
				handlerFunc(c)
			}
		}

		l.Debug("[OPA] OPA is enabled for endpoint ", remote.Endpoint)

		return func(c *gin.Context) {

			res, err := conf.checkPermission(c.Request)
			if err != nil {
				l.Error("[OPA] Error checking permission ", err)
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			if !res {
				l.Error("[OPA] Permission denied")
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "Permission Denied"})
				return
			}

			handlerFunc(c)
		}
	}
}

func (x *xtraConfig) buildRequest(r *http.Request) *Request {
	if r == nil || r.Method == "" || r.URL == nil {
		return nil
	}
	req := &Request{
		Input: Input{
			Method: r.Method,
			Path:   strings.Split(strings.Trim(r.URL.Path, "/"), "/"),
		},
	}

	if x.PayloadMap != nil {
		token := ""
		tHeaders := ""
		tPayload := ""

		bodyRead := false
		tokenRead := false
		var raw []byte

		data := make(map[string]interface{})
		for k, v := range x.PayloadMap {
			if !strings.Contains(v, ".") {
				data[k] = v
				continue
			}
			parts := strings.Split(v, ".")

			switch strings.ToLower(parts[0]) {
			case "jwt":
				if !tokenRead {
					tokenRead = true
					if r.Header == nil {
						continue
					}
					if auth, ok := r.Header[authHeader]; ok {
						token = auth[0]
					}

					if token != "" {
						token = strings.TrimPrefix(token, "Bearer ")
						//Just in case using lower case bearer
						token = strings.TrimPrefix(token, "bearer ")
					}

					tparts := strings.Split(token, ".")
					if len(tparts) != 3 {
						token = ""
						continue
					} else {
						if b, err := base64.RawStdEncoding.DecodeString(tparts[0]); err == nil {
							tHeaders = string(b)
						}

						if b, err := base64.RawStdEncoding.DecodeString(tparts[1]); err == nil {
							tPayload = string(b)
						}
					}
				}
				if tHeaders == "" || tPayload == "" {
					continue
				}
				if parts[1] == "raw" {
					data[k] = token
					continue
				}
				if val := getJWTValue(tHeaders, tPayload, parts[1:]); val != "" {
					data[k] = val
				}
			case "body":
				if !bodyRead {
					bodyRead = true
					if r.Body == nil {
						continue
					}
					raw, _ = ioutil.ReadAll(r.Body)
					if raw == nil {
						continue
					}
					r.Body = ioutil.NopCloser(bytes.NewReader(raw))
				}
				if raw == nil {
					continue
				}
				if parts[1] == "raw" {
					var d map[string]interface{}
					if err := json.Unmarshal(raw, &d); err == nil {
						data[k] = d
					}
					continue
				}

				if val := getBodyValue(string(raw), parts[1:]); val != "" {
					data[k] = val
				}

			case "header":
				if r.Header == nil {
					continue
				}
				if parts[1] == "raw" {
					data[k] = r.Header
					continue
				}

				if val := r.Header.Get(parts[1]); val != "" {
					data[k] = val
				}
			case "query":
				if r.URL.Query() == nil {
					continue
				}
				if parts[1] == "raw" {
					data[k] = r.URL.Query()
					continue
				}

				if val := r.URL.Query().Get(parts[1]); val != "" {
					data[k] = val
				}

			default:
				continue
			}
		}

		req.Input.Payload = data
	}

	return req
}

func (x *xtraConfig) checkPermission(r *http.Request) (bool, error) {

	req := x.buildRequest(r)
	if req == nil {
		return false, errors.New("Fail to build input request")
	}

	return x.Service.Evaluate(x.PackageName, x.Directive, req)

}

func getBodyValue(body string, path []string) interface{} {
	if val := gjson.Get(body, strings.Join(path, ".")); val.Exists() {
		return val.Value()
	}
	return nil
}

func getJWTValue(header string, payload string, path []string) interface{} {
	var data string
	switch strings.ToLower(path[0]) {
	case "header":
		data = header
	case "payload":
		data = payload
	default:
		return nil
	}

	if len(path) == 1 {
		var d map[string]interface{}
		return json.Unmarshal([]byte(data), &d)
	}
	if val := gjson.Get(data, strings.Join(path[1:], ".")); val.Exists() {
		return val.Value()
	}
	return nil
}
