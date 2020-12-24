package opa

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	mc := NewMemoryCache(0)

	mc.Set(hash("test1"), true)
	rsp, ok := mc.Get(hash("test1"))

	assert.True(t, ok)
	assert.True(t, rsp.(bool))

	rsp, ok = mc.Get(hash("test2"))
	assert.False(t, ok)
	assert.Nil(t, rsp)

	mc.Set(hash("test2"), false)
	rsp, ok = mc.Get(hash("test2"))
	assert.True(t, ok)
	assert.False(t, rsp.(bool))
}

func TestCacheExpiration(t *testing.T) {
	mc := NewMemoryCache(time.Duration(10 * time.Millisecond))

	mc.Set(hash("test1"), true)
	rsp, ok := mc.Get(hash("test1"))
	assert.True(t, ok)
	assert.True(t, rsp.(bool))

	time.Sleep(time.Duration(20 * time.Millisecond))
	rsp, ok = mc.Get(hash("test1"))
	assert.False(t, ok)
	assert.Nil(t, rsp)

}

func hash(s string) [32]byte {
	return sha256.Sum256([]byte(s))
}
