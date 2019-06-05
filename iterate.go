package properties

// Map gets the map of properties
func (p Doc) Map() map[string]string {
	m := make(map[string]string)
	p.Foreach(func(v, k string) bool {
		m[k] = v
		return true
	})

	return m
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
