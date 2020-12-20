// Package shared provides space to share variables.
package shared

import (
	"strings"
	"sync"
)

// Vars keeps values of named variables.
type Vars struct {
	// VarPrefix determines which cell values should be collected as vars and replaced with values in usages.
	// Default is '$', e.g. $id1 would be treated as variable.
	VarPrefix string

	mu    sync.Mutex
	vars  map[string]interface{}
	onSet []func(key string, val interface{})
}

// Reset removes all variables.
func (v *Vars) Reset() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.vars = make(map[string]interface{})
	v.onSet = nil
}

// IsVar checks if string looks like a variable name.
func (v *Vars) IsVar(s string) bool {
	varPrefix := v.VarPrefix
	if varPrefix == "" {
		varPrefix = "$"
	}

	return strings.HasPrefix(s, varPrefix)
}

// Get returns variable value if is exists.
func (v *Vars) Get(s string) (interface{}, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	val, found := v.vars[s]

	return val, found
}

// Set sets variable by name.
func (v *Vars) Set(key string, val interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.vars == nil {
		v.vars = make(map[string]interface{})
	}

	v.vars[key] = val

	for _, f := range v.onSet {
		f(key, val)
	}
}

// OnSet adds callback to invoke when variable is set.
//
// All callbacks are removed on Reset.
func (v *Vars) OnSet(f func(key string, val interface{})) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.onSet = append(v.onSet, f)
}

// GetAll returns all variables with values.
func (v *Vars) GetAll() map[string]interface{} {
	v.mu.Lock()
	defer v.mu.Unlock()

	res := make(map[string]interface{}, len(v.vars))
	for k, val := range v.vars {
		res[k] = val
	}

	return res
}
