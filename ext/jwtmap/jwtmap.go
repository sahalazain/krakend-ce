package jwtmap

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	"github.com/tidwall/sjson"
)

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

		l.Debug("[JWTMap] JWTMap is enabled for endpoint ", remote.Endpoint)

		return func(c *gin.Context) {

			if err := extractClaim(conf, c.Request); err != nil {
				l.Error("[JWTMap] Error extracing jwt ", err)
			}

			handlerFunc(c)
		}
	}
}

func extractClaim(config *xtraConfig, r *http.Request) error {
	if config == nil {
		return errors.New("Empty config, skip")
	}

	token := ""
	auth, ok := r.Header[authHeader]
	if ok {
		token = auth[0]
	}

	if token == "" {
		return errors.New("Token is empty, skip extracting ")
	}

	if len(token) < 7 {
		return errors.New("Invalid auth header length")
	}

	token = token[7:]

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("Token is malformed , skip processing")
	}

	tHeaders := ""
	tPayload := ""

	if b, err := base64.RawStdEncoding.DecodeString(parts[0]); err == nil {
		tHeaders = string(b)
	}

	if b, err := base64.RawStdEncoding.DecodeString(parts[1]); err == nil {
		tPayload = string(b)
	}

	if tHeaders == "" || tPayload == "" {
		return errors.New("Token content is not valid")
	}

	for k, v := range config.JWTMap {
		val := getJWTValue(tHeaders, tPayload, strings.Split(v, "."))
		injectResult(k, fmt.Sprintf("%v", val), r)
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

func injectResult(path, val string, r *http.Request) error {
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return errors.New("Invalid result path")
	}

	switch strings.ToLower(parts[0]) {
	case "header":
		r.Header.Set(parts[1], val)
		return nil
	case "body":
		raw, _ := ioutil.ReadAll(r.Body)
		if raw == nil {
			return errors.New("Unable to read request body")
		}
		res, err := sjson.Set(string(raw), strings.Join(parts[1:], "."), val)
		if err != nil {
			return err
		}
		bres := []byte(res)
		r.Body = ioutil.NopCloser(bytes.NewReader(bres))
		r.Header.Set("Content-Length", fmt.Sprintf("%v", len(bres)))
		return nil
	case "query":
		uv := r.URL.Query()
		uv.Add(parts[1], val)
		r.URL.RawQuery = uv.Encode()
		return nil
	default:
		return errors.New("Invalid result path")
	}
}
