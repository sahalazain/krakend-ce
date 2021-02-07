package keyauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChildLookup(t *testing.T) {
	obj := map[string]interface{}{
		"result": map[string]interface{}{
			"name": "sahal",
		},
	}

	name, ok := lookup("result.name", obj)
	assert.True(t, ok)
	assert.Equal(t, "sahal", name)
}

func TestSimpleLookup(t *testing.T) {
	obj := map[string]interface{}{
		"result": "sahal",
	}

	name, ok := lookup("result", obj)
	assert.True(t, ok)
	assert.Equal(t, "sahal", name)
}
