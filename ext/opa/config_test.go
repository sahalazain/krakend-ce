package opa

import (
	"testing"

	"github.com/devopsfaith/krakend/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigInvalidParse(t *testing.T) {

	assert.Nil(t, configGetter(config.ExtraConfig{
		"opa": map[string]interface{}{
			"service_address": "http://localhost:8080",
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: "opa",
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"package_name": "opa.test",
		},
	}), "Should nil")
}

func TestConfigParse(t *testing.T) {
	xtra := config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
			"package_name":    "opa.test",
		},
	}

	cfg := configGetter(xtra)

	assert.NotNil(t, cfg, "Should not nil")
	assert.NotNil(t, cfg.Service)
	assert.Equal(t, cfg.BasePath, basePath, "Should be default")
	assert.Equal(t, cfg.Directive, "allow", "Should be allow")
	assert.Equal(t, cfg.CacheDuration, defaultCacheDuration, "Should be default")
}

func TestConfigCustomParse(t *testing.T) {
	xtra := config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
			"package_name":    "opa.test",
			"directive":       "access",
			"base_path":       "/v2/data",
			"cache_duration":  10,
		},
	}

	cfg := configGetter(xtra)

	assert.NotNil(t, cfg, "Should not nil")
	assert.Equal(t, cfg.BasePath, "/v2/data", "Should not default")
	assert.Equal(t, cfg.Directive, "access", "Should be acccess")
	assert.Equal(t, cfg.CacheDuration, 10, "Should not default")
}

func TestConfigPayloadParse(t *testing.T) {
	xtra := config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
			"package_name":    "opa.test",
			"payload": map[string]interface{}{
				"username": "jwt.claim.username",
				"data":     "body.data",
				"header":   "header.version",
				"version":  2,
			},
		},
	}

	cfg := configGetter(xtra)
	assert.NotNil(t, cfg, "Should not nil")

	assert.Equal(t, cfg.PayloadMap["data"], "body.data", "Should be body.data")
	assert.Equal(t, cfg.PayloadMap["version"], "", "Should be nil")
}
