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

	expect(t, "注释之后保存", "#This is a comment for a\na=aaa\n" == buf)
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
	if nil != err {
		t.Error("加载失败")
		return
	}

	v := ""

	v = p.Str("a")
	if "aa" != v {
		t.Error("Get string failed")
		return
	}

	v = p.Str("b")
	if "bbb" != v {
		t.Error("Get string failed")
		return
	}

	v = p.Str("Z")
	if "" != v {
		t.Error("Get string failed")
		return
	}

	v = p.Str("c ccc")
	if "cccc" != v {
		t.Error("Get string failed")
		return
	}

	v = p.Str("dd")
	if "" != v {
		t.Error("Get string failed")
		return
	}

	v = p.Str("ee")
	if "r-rt rr" != v {
		t.Error("Get string failed")
		return
	}
}

func Test_LoadFromFile(t *testing.T) {
	_, err := LoadFile("notexists.properties")
	assert.NotNil(t, err)

	doc, err := LoadFile("test1.properties")
	if nil != err {
		t.Error("加载失败")
		return
	}

	fmt.Println(doc.Str("key"))
}
