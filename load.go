package properties

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"
)

// LoadFile creates the properties document from a file or a stream.
func LoadFile(file string) (doc *Doc, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return Load(f)
}

// LoadMap creates a new Properties struct from a string map.
// copy from https://github.com/magiconair/properties
func LoadMap(m map[string]string) (doc *Doc, err error) {
	p := New()

	for k, v := range m {
		p.Set(k, v)
	}

	return p, nil
}

// LoadString creates the properties document from a string.
func LoadString(s string) (doc *Doc, err error) {
	return Load(strings.NewReader(s))
}

// LoadBytes creates the properties document from a string.
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
		l := scanner.Bytes()

		//  遇到空行
		if len(l) == 0 {
			doc.lines.PushBack(&line{typo: ' ', value: string("")})
			continue
		}

		//  找到第一个非空白字符
		pos := bytes.IndexFunc(l, func(r rune) bool { return !unicode.IsSpace(r) })

		//  遇到空白行
		if pos == -1 {
			doc.lines.PushBack(&line{typo: ' ', value: string("")})
			continue
		}

		//  遇到注释行
		if isComment(l[pos]) {
			doc.lines.PushBack(&line{typo: l[pos], value: string(l)})
			continue
		}

		//  找到第一个等号的位置
		end := bytes.IndexFunc(l[pos+1:], func(r rune) bool { return r == '=' || r == ':' })

		var (
			typo       byte = '=' //  没有=，说明该配置项只有key
			key, value []byte
		)

		if end == -1 {
			key = bytes.TrimRightFunc(l[pos:], unicode.IsSpace)
		} else {
			key = bytes.TrimRightFunc(l[pos:pos+1+end], unicode.IsSpace)
			value = bytes.TrimSpace(l[pos+1+end+1:])
			typo = l[pos+1+end]
		}

		elem := &line{typo: typo, key: string(key), value: string(value)}
		doc.props[string(key)] = doc.lines.PushBack(elem)
	}

	return doc, scanner.Err()
}
