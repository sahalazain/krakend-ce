package jwtmap

import (
	"testing"

	"github.com/devopsfaith/krakend/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.Nil(t, configGetter(config.ExtraConfig{
		"jwt_map": map[string]interface{}{
			"name": "sahal",
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"jwt_map": map[string]interface{}{
				"name": "sahal",
			},
		},
	}), "Should nil")

	assert.Nil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"jwt_map": map[string]interface{}{
				"body.name": "sahal",
			},
		},
	}), "Should nil")

	assert.NotNil(t, configGetter(config.ExtraConfig{
		namespace: map[string]interface{}{
			"jwt_map": map[string]interface{}{
				"body.name": "payload.name",
			},
		},
	}), "Should not nil")
}
