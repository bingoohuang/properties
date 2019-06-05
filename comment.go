package properties

import (
	"bufio"
	"strings"
)

// Comment appends comments for the special line.
//
// Return false if the special line is not exist.
func (p *Doc) Comment(key, comments string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	//  如果所有注释为空
	if comments == "" {
		p.lines.InsertBefore(&line{typo: '#', value: "#"}, e)
		return true
	}

	//  创建一个新的Scanner
	scanner := bufio.NewScanner(strings.NewReader(comments))
	for scanner.Scan() {
		p.lines.InsertBefore(&line{typo: '#', value: "#" + scanner.Text()}, e)
	}

	return true
}

// Uncomment removes all of the comments for the special line.
//
// Return false if the special line is not exist.
func (p *Doc) Uncomment(key string) bool {
	e, ok := p.props[key]
	if !ok {
		return false
	}

	for i := e.Prev(); nil != i; {
		del := i
		i = i.Prev()

		typo := del.Value.(*line).typo
		if isProperty(typo) || typo == ' ' {
			break
		}

		p.lines.Remove(del)
	}

	return true
}
