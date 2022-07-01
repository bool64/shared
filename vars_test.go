package shared_test

import (
	"context"
	"sync"
	"testing"

	"github.com/bool64/shared"
	"github.com/stretchr/testify/assert"
)

func TestVars_GetAll(t *testing.T) {
	v := shared.Vars{}
	cnt := 0

	v.OnSet(func(key string, val interface{}) {
		assert.Equal(t, "k", key)
		assert.Equal(t, "v", val)

		cnt++
	})

	v.Set("k", "v")
	val, found := v.Get("k")
	assert.True(t, found)
	assert.Equal(t, "v", val)

	assert.Equal(t, map[string]interface{}{"k": "v"}, v.GetAll())
	assert.True(t, v.IsVar("$var"))
	assert.False(t, v.IsVar("var"))
	assert.Equal(t, 1, cnt)

	v.Reset()
	assert.Equal(t, map[string]interface{}{}, v.GetAll())
}

func TestVars_Fork(t *testing.T) {
	v := shared.Vars{}
	v.Set("k", "v")

	wg := sync.WaitGroup{}

	for i := 0; i < 50; i++ {
		i := i

		wg.Add(1)

		go func() {
			defer wg.Done()

			ctx, vi := v.Fork(context.Background())

			assert.Equal(t, map[string]interface{}{"k": "v"}, vi.GetAll())
			vi.Set("ki", i)
			assert.Equal(t, map[string]interface{}{"k": "v", "ki": i}, vi.GetAll())
			vi.Set("k", i)
			assert.Equal(t, map[string]interface{}{"k": i, "ki": i}, vi.GetAll())

			// Forking with already instrumented context is a no op.
			ctx2, vi2 := v.Fork(ctx)
			assert.Equal(t, ctx, ctx2)
			assert.Equal(t, vi, vi2)
		}()
	}

	assert.Equal(t, map[string]interface{}{"k": "v"}, v.GetAll())
}

func TestVarToContext(t *testing.T) {
	ctx := context.Background()
	ctx = shared.VarToContext(ctx, "$foo", "bar")

	assert.Equal(t, map[string]interface{}{"$foo": "bar"}, shared.VarsFromContext(ctx))

	var nilParent *shared.Vars

	npCtx, vv := nilParent.Fork(ctx)
	assert.Equal(t, map[string]interface{}{"$foo": "bar"}, shared.VarsFromContext(npCtx))
	assert.Equal(t, map[string]interface{}{"$foo": "bar"}, vv.GetAll())

	parent := shared.Vars{}
	parent.Set("$baz", "qux")

	pCtx, vv := parent.Fork(ctx)
	assert.Equal(t, map[string]interface{}{"$foo": "bar", "$baz": "qux"}, shared.VarsFromContext(pCtx))
	assert.Equal(t, map[string]interface{}{"$foo": "bar", "$baz": "qux"}, vv.GetAll())

	ctx = shared.VarToContext(pCtx, "$quux", true)
	assert.Equal(t, map[string]interface{}{"$foo": "bar", "$baz": "qux", "$quux": true}, shared.VarsFromContext(ctx))
}
