package canlog

import (
	"context"
	"sync"
)

// Logger methods
type Logger interface {
	Log(msg string, fields map[string]interface{})
}

// DefaultLogger is used by default (singleton)
var DefaultLogger Logger

type ctxKey struct{}

type ctxValue struct {
	fields map[string]interface{}
	*sync.Mutex
}

// NewContext return a new ctx with a logger
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, ctxValue{
		map[string]interface{}{},
		&sync.Mutex{},
	})
}

// find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func find(slice []interface{}, val interface{}) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// AppendErr to a slice in the canonical logger
func AppendErr(ctx context.Context, value interface{}) {
	ctxVal := ctx.Value(ctxKey{}).(ctxValue)
	ctxVal.Lock()

	current, ok := ctxVal.fields["errid"].([]interface{})
	if !ok {
		current = []interface{}{}
		ctxVal.fields["errid"] = current
	}

	// Prevent dupes
	_, exists := find(current, value)
	if exists {
		ctxVal.Unlock()
		return
	}
	current = append(current, value)
	ctxVal.fields["errid"] = current

	ctxVal.Unlock()
}

// Append to a slice in the canonical logger
func Append(ctx context.Context, key string, value interface{}) {
	ctxVal := ctx.Value(ctxKey{}).(ctxValue)
	ctxVal.Lock()

	current, ok := ctxVal.fields[key].([]interface{})
	if !ok {
		current = []interface{}{}
		ctxVal.fields[key] = current
	}

	// Prevent dupes
	_, exists := find(current, value)
	if exists {
		ctxVal.Unlock()
		return
	}
	current = append(current, value)
	ctxVal.fields[key] = current

	ctxVal.Unlock()
}

const CtxErrKey = "err"

// SetErr a value for logging
func SetErr(ctx context.Context, err error) {
	ctxVal := ctx.Value(ctxKey{}).(ctxValue)
	ctxVal.Lock()
	ctxVal.fields[CtxErrKey] = err.Error()
	ctxVal.Unlock()
}

// Set a value for logging
func Set(ctx context.Context, key string, value interface{}) {
	ctxVal := ctx.Value(ctxKey{}).(ctxValue)
	ctxVal.Lock()
	ctxVal.fields[key] = value
	ctxVal.Unlock()
}

// Log full canonical log line to logger
func Log(ctx context.Context, msg string) {
	ctxVal := ctx.Value(ctxKey{}).(ctxValue)
	ctxVal.Lock()
	DefaultLogger.Log(msg, ctxVal.fields)
	ctxVal.fields = map[string]interface{}{}
	ctxVal.Unlock()
}
