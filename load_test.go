package properties

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Save(t *testing.T) {
	doc := New()
	doc.Set("a", "aaa")
	doc.Comment("a", "This is a comment for a")

	buf, err := doc.Export()
	assert.Nil(t, err)

	expect(t, "注释之后保存", buf == "#This is a comment for a\na=aaa\n")
	expect(t, "注释之后保存", buf == doc.String())
}

func Test_Load(t *testing.T) {
	s := `
    a=aa
    b=bbb
    c ccc = cccc
    dd
    # commment1
    !comment2
    
    ee: r-rt rr
    `

	p, err := LoadString(s)
	if err != nil {
		t.Error("加载失败")
		return
	}

	v := ""

	v = p.Str("a")
	if v != "aa" {
		t.Error("Get string failed")
		return
	}

	v = p.Str("b")
	if v != "bbb" {
		t.Error("Get string failed")
		return
	}

	v = p.Str("Z")
	if v != "" {
		t.Error("Get string failed")
		return
	}

	v = p.Str("c ccc")
	if v != "cccc" {
		t.Error("Get string failed")
		return
	}

	v = p.Str("dd")
	if v != "" {
		t.Error("Get string failed")
		return
	}

	v = p.Str("ee")
	if v != "r-rt rr" {
		t.Error("Get string failed")
		return
	}
}

func Test_LoadFromFile(t *testing.T) {
	_, err := LoadFile("notexists.properties")
	assert.NotNil(t, err)

	doc, err := LoadFile("load_test.properties")
	if err != nil {
		t.Error("加载失败")
		return
	}

	fmt.Println(doc.Str("key"))
}
