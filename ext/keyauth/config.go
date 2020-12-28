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
	defaultHeaderName    = "X-KeyID"
)

type xtraConfig struct {
	ServiceAddress string
	BasePath       string
	KeyPath        string
	IDHeaderName   string
	CacheDuration  int
	Service        service.KeyAuth
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
		BasePath:      basePath,
		IDHeaderName:  defaultHeaderName,
	}

	if sa, ok := tmp["service_address"].(string); ok {
		conf.ServiceAddress = sa
	} else {
		return nil
	}

	if kp, ok := tmp["key_path"].(string); ok {
		if !strings.Contains(kp, ".") {
			return nil
		}
		conf.KeyPath = kp
	} else {
		return nil
	}

	if cd, ok := tmp["cache_duration"]; ok {
		fmt.Println("Change cache duration")
		if cdi, err := strconv.Atoi(fmt.Sprintf("%v", cd)); err == nil {
			conf.CacheDuration = cdi
		}
	}

	if bp, ok := tmp["base_path"].(string); ok {
		conf.BasePath = bp
	}

	if hn, ok := tmp["header_name"].(string); ok {
		conf.IDHeaderName = hn
	}

	conf.Service = service.NewHTTPKeyAuth(conf.ServiceAddress, conf.BasePath, conf.CacheDuration)

	return &conf
}
