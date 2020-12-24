package opa

import (
	"fmt"
	"strconv"

	"github.com/devopsfaith/krakend-ce/ext/service"
	"github.com/devopsfaith/krakend/config"
)

const (
	basePath             = "/v1/data/"
	namespace            = "github_com/sahalzain/krakend-opa"
	authHeader           = "Authorization"
	defaultCacheDuration = 24 * 3600
)

type xtraConfig struct {
	ServiceAddress string
	PackageName    string
	BasePath       string
	Directive      string
	PayloadMap     map[string]string
	CacheDuration  int
	Service        service.Policy
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
		Directive:     "allow",
		BasePath:      basePath,
		CacheDuration: defaultCacheDuration,
	}

	if sa, ok := tmp["service_address"].(string); ok {
		conf.ServiceAddress = sa
	} else {
		return nil
	}

	if pkg, ok := tmp["package_name"].(string); ok {
		conf.PackageName = pkg
	} else {
		return nil
	}

	if dr, ok := tmp["directive"].(string); ok {
		conf.Directive = dr
	}

	if bp, ok := tmp["base_path"].(string); ok {
		conf.BasePath = bp
	}

	if cd, ok := tmp["cache_duration"]; ok {
		fmt.Println("Change cache duration")
		if cdi, err := strconv.Atoi(fmt.Sprintf("%v", cd)); err == nil {
			conf.CacheDuration = cdi
		}
	}
	if pm, ok := tmp["payload"].(map[string]interface{}); ok {
		tmp := make(map[string]string)
		for k, v := range pm {
			switch vt := v.(type) {
			case string:
				tmp[k] = vt
			default:
				continue
			}
		}
		conf.PayloadMap = tmp
	}

	conf.Service = service.NewHTTPOPA(conf.ServiceAddress, conf.BasePath, conf.CacheDuration)

	return &conf
}
