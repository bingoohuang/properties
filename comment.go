package properties

import (
	"bufio"
	"strings"
)

// Comment appends comments for the special item.
//
// Return false if the special item is not exist.
func (p *Doc) Comment(key, comments string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	//  如果所有注释为空
	if comments == "" {
		p.items.InsertBefore(&item{typo: '#', value: "#"}, e)
		return true
	}

	//  创建一个新的Scanner
	scanner := bufio.NewScanner(strings.NewReader(comments))
	for scanner.Scan() {
		p.items.InsertBefore(&item{typo: '#', value: "#" + scanner.Text()}, e)
	}

	return true
}

// Uncomment removes all of the comments for the special item.
//
// Return false if the special item is not exist.
func (p *Doc) Uncomment(key string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	for i := e.Prev(); nil != i; {
		del := i
		i = i.Prev()

		typo := del.Value.(*item).typo
		if typo == '=' || typo == ':' || typo == ' ' {
			break
		}

		p.items.Remove(del)
	}

	return true
}
