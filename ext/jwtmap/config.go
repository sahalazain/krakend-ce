package jwtmap

import (
	"fmt"
	"strings"

	"github.com/devopsfaith/krakend/config"
)

const (
	namespace  = "github_com/sahalzain/krakend-jwtmap"
	authHeader = "Authorization"
)

type xtraConfig struct {
	JWTMap map[string]string
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
		JWTMap: make(map[string]string),
	}

	if rm, ok := tmp["jwt_map"]; ok {
		if rmap, ok := rm.(map[string]interface{}); ok {
			for k, v := range rmap {
				if !strings.Contains(fmt.Sprintf("%v", v), ".") {
					continue
				}

				if !strings.Contains(fmt.Sprintf("%v", k), ".") {
					continue
				}

				conf.JWTMap[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	if len(conf.JWTMap) == 0 {
		return nil
	}

	return &conf
}
