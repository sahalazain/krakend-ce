package keyauth

import (
	"testing"

	"github.com/devopsfaith/krakend/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigInvalidParse(t *testing.T) {

	assert.Nil(t, configGetter(config.ExtraConfig{
		"keyauth": map[string]interface{}{
			"service_address": "http://localhost:8080",
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: "keyauth",
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"key_path": "body.key_api",
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
			"key_path":        "body",
		},
	}))
}

func TestConfigParse(t *testing.T) {
	xtra := config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
			"request_map": map[string]interface{}{
				"key": "body.key_api",
			},
		},
	}

	cfg := configGetter(xtra)

	assert.NotNil(t, cfg, "Should not nil")
	assert.NotNil(t, cfg.Service)
	assert.Equal(t, cfg.BasePath, basePath, "Should be default")
	assert.Equal(t, cfg.CacheDuration, defaultCacheDuration, "Should be default")
}

func TestConfigCustomParse(t *testing.T) {
	xtra := config.ExtraConfig{
		namespace: map[string]interface{}{
			"service_address": "http://localhost:8080",
			"package_name":    "opa.test",
			"request_map": map[string]interface{}{
				"key": "body.key_api",
			},
			"base_path":      "/v2/auth/key",
			"cache_duration": 10,
			"response_map": map[string]interface{}{
				"header.X-PartnerID": "result.partner",
			},
		},
	}

	cfg := configGetter(xtra)

	assert.NotNil(t, cfg, "Should not nil")
	assert.Equal(t, cfg.BasePath, "/v2/auth/key", "Should not default")
	assert.Equal(t, cfg.CacheDuration, 10, "Should not default")
	assert.Equal(t, 1, len(cfg.ResponseMap))
}
