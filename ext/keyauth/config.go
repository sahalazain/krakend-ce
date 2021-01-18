package keyauth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/devopsfaith/krakend-ce/ext/service"
	"github.com/devopsfaith/krakend/config"
)

const (
	namespace            = "github_com/sahalzain/krakend-keyauth"
	basePath             = "/v1/auth/key"
	defaultCacheDuration = 24 * 3600
	defaultResultPath    = "header.X-KeyID"
	defaultResponsePath  = "result"
)

type xtraConfig struct {
	ServiceAddress string
	BasePath       string
	IDResultPath   string
	CacheDuration  int
	CacheSize      int
	Service        service.KeyAuth
	RequestMap     map[string]string
	ResponsePath   string
}

func configGetter(cfg config.ExtraConfig) *xtraConfig {
	v, ok := cfg[namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	conf := xtraConfig{
		CacheDuration: defaultCacheDuration,
		CacheSize:     0,
		BasePath:      basePath,
		IDResultPath:  defaultResultPath,
		RequestMap:    make(map[string]string),
	}

	if sa, ok := tmp["service_address"].(string); ok {
		conf.ServiceAddress = sa
	} else {
		return nil
	}

	if rm, ok := tmp["request_map"]; ok {
		if rmap, ok := rm.(map[string]interface{}); ok {
			for k, v := range rmap {
				if !strings.Contains(fmt.Sprintf("%v", v), ".") {
					continue
				}
				conf.RequestMap[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	if len(conf.RequestMap) == 0 {
		return nil
	}

	if cd, ok := tmp["cache_duration"]; ok {
		if cdi, err := strconv.Atoi(fmt.Sprintf("%v", cd)); err == nil {
			conf.CacheDuration = cdi
		}
	}

	if cs, ok := tmp["cache_size"]; ok {
		if csi, err := strconv.Atoi(fmt.Sprintf("%v", cs)); err == nil {
			conf.CacheSize = csi
		}
	}

	if rp, ok := tmp["response_path"].(string); ok {
		conf.ResponsePath = rp
	}

	if bp, ok := tmp["base_path"].(string); ok {
		conf.BasePath = bp
	}

	if hn, ok := tmp["result_path"].(string); ok {
		if strings.Contains(hn, ".") {
			conf.IDResultPath = hn
		}
	}

	conf.Service = service.NewHTTPKeyAuth(conf.ServiceAddress, conf.BasePath, conf.ResponsePath, conf.CacheDuration, conf.CacheSize)

	return &conf
}
