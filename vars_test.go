package shared_test

import (
	"testing"

	"github.com/bool64/shared"
	"github.com/stretchr/testify/assert"
)

func TestVars_GetAll(t *testing.T) {
	v := shared.Vars{}

	v.Set("k", "v")
	val, found := v.Get("k")
	assert.True(t, found)
	assert.Equal(t, "v", val)

	assert.Equal(t, map[string]interface{}{"k": "v"}, v.GetAll())
	assert.True(t, v.IsVar("$var"))
	assert.False(t, v.IsVar("var"))

	v.Reset()
	assert.Equal(t, map[string]interface{}{}, v.GetAll())
}
