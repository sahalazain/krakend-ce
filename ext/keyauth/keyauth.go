package keyauth

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendgin "github.com/devopsfaith/krakend/router/gin"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//Request KeyAuth request model
type Request struct {
	Key string `json:"key,omitempty" mapstructure:"key"`
	ID  string `json:"id,omitempty" mapstructure:"id"`
}

//Hash calculate request hash
func (r *Request) Hash() [32]byte {
	val := fmt.Sprintf("%v", r)
	return sha256.Sum256([]byte(val))
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

		l.Debug("[KeyAuth] KeyAuth is enabled for endpoint ", remote.Endpoint)

		return func(c *gin.Context) {
			res, err := conf.validateKey(c.Request)
			if err != nil {
				l.Error("[KeyAuth] Error validating key ", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
				return
			}

			if !res {
				l.Error("[KeyAuth] Invalid Key API")
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid KEY API"})
				return
			}

			handlerFunc(c)
		}
	}
}

func (x *xtraConfig) validateKey(r *http.Request) (bool, error) {
	key, err := x.extractKey(r)
	if err != nil {
		return false, err
	}

	req := &Request{
		Key: key,
	}

	id, err := x.Service.Validate(req)
	if err != nil {
		return false, err
	}

	if id == "" {
		return false, nil
	}

	r.Header.Set(x.IDHeaderName, id)

	return true, nil
}

func (x *xtraConfig) extractKey(r *http.Request) (string, error) {
	parts := strings.Split(x.KeyPath, ".")
	switch strings.ToLower(parts[0]) {
	case "body":
		if r.Body == nil {
			return "", errors.New("Empty request body")
		}

		raw, _ := ioutil.ReadAll(r.Body)
		if raw == nil {
			return "", errors.New("Unable to read request body")
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(raw))

		if val := gjson.Get(string(raw), strings.Join(parts[1:], ".")); val.Exists() {
			return val.String(), nil
		}

		return "", errors.New("API Key on body not found")

	case "header":
		if r.Header == nil {
			return "", errors.New("Empty header")
		}

		if val := r.Header.Get(parts[1]); val != "" {
			return val, nil
		}

		return "", errors.New("API Key on header not found")
	case "query":
		if r.URL.Query() == nil {
			return "", errors.New("Empty url query")
		}

		if val := r.URL.Query().Get(parts[1]); val != "" {
			return val, nil
		}

		return "", errors.New("API Key on query not found")
	default:
		return "", errors.New("Invalid path")
	}
}
