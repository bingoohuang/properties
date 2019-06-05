// Package properties is used to read or write or modify the properties document.
package properties

import (
	"container/list"
	"strconv"
)

// New  creates a new and empty properties document.
//
// It's used to generate a new document.
func New() *Doc {
	return &Doc{
		items: list.New(),
		props: make(map[string]*list.Element),
	}
}

// Get retrieves the value from Doc.
//
// If the item is not exist, the exist is false.
func (p Doc) Get(key string) (value string, exist bool) {
	if e, ok := p.props[key]; ok {
		return e.Value.(*item).value, ok
	}

	return "", false
}

// MustGet returns the expanded value for the given key if exists or
// panics otherwise.
func (p Doc) MustGet(key string) (value string) {
	if e, ok := p.props[key]; ok {
		return e.Value.(*item).value
	}

	panic(key + " not found")
}

// Set updates the value of the item of the key.
//
// Create a new item if the item of the key is not exist.
func (p *Doc) Set(key, value string) {
	if e, ok := p.props[key]; ok {
		e.Value.(*item).value = value
	} else {
		p.props[key] = p.items.PushBack(&item{typo: '=', key: key, value: value})
	}
}

// Del deletes the exist item.
//
// If the item is not exist, return false.
func (p *Doc) Del(key string) bool {
	if e, ok := p.props[key]; ok {
		p.Uncomment(key)
		p.items.Remove(e)
		delete(p.props, key)
		return true
	}

	return false
}

// Accept traverses every item of the document, include comment.
//
// The typo parameter special the item type.
// If typo is '#' or '!' means current item is a comment.
// If typo is ' ' means current item is a empty or a space line.
// If typo is '=' or ':' means current item is a key-value pair.
// The traverse will be terminated if f return false.
func (p Doc) Accept(f func(typo byte, value, key string) bool) {
	for e := p.items.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*item)
		if continues := f(elem.typo, elem.value, elem.key); !continues {
			return
		}
	}
}

// Foreach traverses all of the key-value pairs in the document.
// The traverse will be terminated if f return false.
func (p Doc) Foreach(f func(value, key string) bool) {
	for e := p.items.Front(); e != nil; e = e.Next() {
		elem := e.Value.(*item)
		if elem.typo == '=' || elem.typo == ':' {
			if continues := f(elem.value, elem.key); !continues {
				return
			}
		}
	}
}

// StringOr retrieves the string value by key.
// If the item is not exist, the def will be returned.
func (p Doc) StringOr(key, def string) string {
	if e, ok := p.props[key]; ok {
		return e.Value.(*item).value
	}

	return def
}

// IntOr retrieves the int value by key.
// If the item is not exist, the def will be returned.
func (p Doc) IntOr(key string, def int) int {
	if e, ok := p.props[key]; ok {
		if v, err := strconv.Atoi(e.Value.(*item).value); err == nil {
			return v
		}
	}

	return def
}

// Int64Or retrieves the int64 value by key.
// If the item is not exist, the def will be returned.
func (p Doc) Int64Or(key string, def int64) int64 {
	if e, ok := p.props[key]; ok {
		if v, err := strconv.ParseInt(e.Value.(*item).value, 10, 64); err == nil {
			return v
		}
	}

	return def
}

// Uint64Or Same as Int64Or, but the return type is uint64.
func (p Doc) Uint64Or(key string, def uint64) uint64 {
	if e, ok := p.props[key]; ok {
		if v, err := strconv.ParseUint(e.Value.(*item).value, 10, 64); err == nil {
			return v
		}
	}

	return def
}

// Float64Or   retrieve the float64 value by key.
// If the item is not exist, the def will be returned.
func (p Doc) Float64Or(key string, def float64) float64 {
	if e, ok := p.props[key]; ok {
		if v, err := strconv.ParseFloat(e.Value.(*item).value, 64); err == nil {
			return v
		}
	}

	return def
}

// BoolOr   retrieve the bool value by key.
// If the item is not exist, the def will be returned.
// This function mapping "1", "t", "T", "true", "TRUE", "True" as true.
// This function mapping "0", "f", "F", "false", "FALSE", "False" as false.
// If the item is not exist of can not map to value of bool,the def will be returned.
func (p Doc) BoolOr(key string, def bool) bool {
	if e, ok := p.props[key]; ok {
		if v, err := strconv.ParseBool(e.Value.(*item).value); err == nil {
			return v
		}
	}

	return def
}

// ObjectOr maps the value of the key to any object.
// The f is the customized mapping function.
// Return def if the item is not exist of f have a error returned.
func (p Doc) ObjectOr(key string, def interface{}, f func(k, v string) (interface{}, error)) interface{} {
	if e, ok := p.props[key]; ok {
		if v, err := f(key, e.Value.(*item).value); err == nil {
			return v
		}
	}

	return def
}

// String same as StringOr but the def is "".
func (p Doc) String(key string) string {
	return p.StringOr(key, "")
}

// Int is same as IntOr but the def is 0 .
func (p Doc) Int(key string) int {
	return p.IntOr(key, 0)
}

// Int64 is same as Int64Or but the def is 0 .
func (p Doc) Int64(key string) int64 {
	return p.Int64Or(key, 0)
}

// Uint64 same as Uint64Or but the def is 0 .
func (p Doc) Uint64(key string) uint64 {
	return p.Uint64Or(key, 0)
}

// Float64 same as Float64Or but the def is 0.0 .
func (p Doc) Float64(key string) float64 {
	return p.Float64Or(key, 0.0)
}

// Bool same as BoolOr but the def is false.
func (p Doc) Bool(key string) bool {
	return p.BoolOr(key, false)
}

// Object is same as ObjectOr but the def is nil.
//
// Notice: If the return value can not be assign to nil, this function will panic/
func (p Doc) Object(key string, f func(k, v string) (interface{}, error)) interface{} {
	return p.ObjectOr(key, nil, f)
}
