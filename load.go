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
