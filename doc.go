// Package properties is used to read or write or modify the properties document.
package properties

import (
	"container/list"
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
	if val, ok := p.Get(key); ok {
		return val
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
