package properties

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func expect(t *testing.T, msg string, result bool) {
	if !result {
		t.Error(msg)
	}
}

func TestDoc_MustGet(t *testing.T) {
	doc := New()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()

	doc.MustGet("key")
}

func Test_New(t *testing.T) {
	doc := New()
	doc.Set("a", "aaa")
	doc.Comment("a", "This is a comment for a")

	buf := bytes.NewBufferString("")
	assert.Nil(t, doc.Save(buf))

	if buf.String() != "#This is a comment for a\na=aaa\n" {
		fmt.Println("Dump failed:[" + buf.String() + "]")
		t.Error("Dump failed")
		return
	}
}

const str = `
	key1=1
	key 2 = 2
	`

func Test_Get(t *testing.T) {
	doc, _ := LoadBytes([]byte(str))

	value1 := doc.MustGet("key1")
	expect(t, "检测Get函数的行为:EXIST", value1 == "1")

	value, exist := doc.Get("NOT-EXIST")
	expect(t, "检测Get函数的行为:NOT-EXIST", !exist)
	expect(t, "检测Get函数的行为:NOT-EXIST", value == "")
}

func Test_Set(t *testing.T) {
	doc, _ := LoadString(str)

	doc.Set("key1", "new-value")
	newValue, _ := doc.Get("key1")
	expect(t, "修改已经存在的项的值", newValue == "new-value")

	doc.Set("NOT-EXIST", "Setup-New-Item")
	newValue, _ = doc.Get("NOT-EXIST")
	expect(t, "修改不存在的项,默认是新增行为", newValue == "Setup-New-Item")
}

func Test_Del(t *testing.T) {
	doc, _ := LoadString(str)

	exist := doc.Del("NOT-EXIST")
	expect(t, "删除不存在的项,需要返回false", !exist)

	exist = doc.Del("key1")
	expect(t, "删除已经存在的项,返回true", exist)
}
