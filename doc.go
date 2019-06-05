// Package properties is used to read or write or modify the properties document.
package properties

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type item struct {
	//  #   注释行
	//  !   注释行
	//  ' ' 空白行或者空行
	//  =   等号分隔的属性行
	//  :   冒号分隔的属性行
	typo  byte   //  行类型
	value string //  行的内容,如果是注释注释引导符也包含在内
	key   string //  如果是属性行这里表示属性的key
}

// Doc The properties document in memory.
type Doc struct {
	items *list.List
	props map[string]*list.Element
}

// New  creates a new and empty properties document.
//
// It's used to generate a new document.
func New() *Doc {
	return &Doc{
		items: list.New(),
		props: make(map[string]*list.Element),
	}
}

// Save saves the doc to file or stream.
func (p Doc) Export() (string, error) {
	buf := bytes.NewBufferString("")
	if err := p.Save(buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Save saves the doc to file or stream.
func (p Doc) Save(writer io.Writer) error {
	var err error

	p.Accept(func(typo byte, value string, key string) bool {
		switch typo {
		case '#', '!', ' ':
			_, err = fmt.Fprintln(writer, value)
		case '=', ':':
			_, err = fmt.Fprintf(writer, "%s%c%s\n", key, typo, value)
		}

		return nil == err
	})

	return err
}

// LoadFile creates the properties document from a file or a stream.
func LoadFile(file string) (doc *Doc, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return Load(f)
}

// LoadString creates the properties document from a string.
func LoadString(s string) (doc *Doc, err error) {
	return Load(strings.NewReader(s))
}

// LoadString creates the properties document from a string.
func LoadBytes(s []byte) (doc *Doc, err error) {
	return Load(bytes.NewReader(s))
}

// Load creates the properties document from a file or a stream.
func Load(reader io.Reader) (doc *Doc, err error) {
	//  创建一个Properties对象
	doc = New()

	//  创建一个扫描器
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		//  逐行读取
		line := scanner.Bytes()

		//  遇到空行
		if 0 == len(line) {
			doc.items.PushBack(&item{typo: ' ', value: string("")})
			continue
		}

		//  找到第一个非空白字符
		pos := bytes.IndexFunc(line, func(r rune) bool { return !unicode.IsSpace(r) })

		//  遇到空白行
		if -1 == pos {
			doc.items.PushBack(&item{typo: ' ', value: string("")})
			continue
		}

		//  遇到注释行
		if line[pos] == '#' || line[pos] == '!' {
			doc.items.PushBack(&item{typo: line[pos], value: string(line)})
			continue
		}

		//  找到第一个等号的位置
		end := bytes.IndexFunc(line[pos+1:], func(r rune) bool { return r == '=' || r == ':' })

		//  没有=，说明该配置项只有key
		var typo byte = '='
		var key []byte
		var value []byte
		if -1 == end {
			key = bytes.TrimRightFunc(line[pos:], unicode.IsSpace)
		} else {
			key = bytes.TrimRightFunc(line[pos:pos+1+end], unicode.IsSpace)
			value = bytes.TrimSpace(line[pos+1+end+1:])
			typo = line[pos+1+end]
		}

		elem := &item{typo: typo, key: string(key), value: string(value)}
		doc.props[string(key)] = doc.items.PushBack(elem)
	}

	return doc, scanner.Err()
}

// Get retrieves the value from PropertiesDocument.
//
// If the item is not exist, the exist is false.
func (p Doc) Get(key string) (value string, exist bool) {
	e, ok := p.props[key]
	if !ok {
		return "", ok
	}

	return e.Value.(*item).value, ok
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
