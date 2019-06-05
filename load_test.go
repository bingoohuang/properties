package properties

import (
	"fmt"
	"testing"
)

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

	v = p.String("a")
	if "aa" != v {
		t.Error("Get string failed")
		return
	}

	v = p.String("b")
	if "bbb" != v {
		t.Error("Get string failed")
		return
	}

	v = p.String("Z")
	if "" != v {
		t.Error("Get string failed")
		return
	}

	v = p.String("c ccc")
	if "cccc" != v {
		t.Error("Get string failed")
		return
	}

	v = p.String("dd")
	if "" != v {
		t.Error("Get string failed")
		return
	}

	v = p.String("ee")
	if "r-rt rr" != v {
		t.Error("Get string failed")
		return
	}
}

func Test_LoadFromFile(t *testing.T) {
	doc, err := LoadFile("test1.properties")
	if nil != err {
		t.Error("加载失败")
		return
	}

	fmt.Println(doc.String("key"))
}
