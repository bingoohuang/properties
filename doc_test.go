package properties

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func expect(t *testing.T, msg string, result bool) {
	if result {
		return
	}

	t.Error(msg)
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

func Test_New(t *testing.T) {
	doc := New()
	doc.Set("a", "aaa")
	doc.Comment("a", "This is a comment for a")

	buf := bytes.NewBufferString("")
	assert.Nil(t, doc.Save(buf))

	if "#This is a comment for a\na=aaa\n" != buf.String() {
		fmt.Println("Dump failed:[" + buf.String() + "]")
		t.Error("Dump failed")
		return
	}
}

func Test_Save(t *testing.T) {
	doc := New()
	doc.Set("a", "aaa")
	doc.Comment("a", "This is a comment for a")

	buf, err := doc.Export()
	assert.Nil(t, err)

	expect(t, "注释之后保存", "#This is a comment for a\na=aaa\n" == buf)
}

const str = `
	key1=1
	key 2 = 2
	`

func Test_Get(t *testing.T) {
	doc, _ := Load(bytes.NewBufferString(str))

	value, exist := doc.Get("key1")
	expect(t, "检测Get函数的行为:EXIST", true == exist)
	expect(t, "检测Get函数的行为:EXIST", value == "1")

	value, exist = doc.Get("NOT-EXIST")
	expect(t, "检测Get函数的行为:NOT-EXIST", false == exist)
	expect(t, "检测Get函数的行为:NOT-EXIST", value == "")
}

func Test_Set(t *testing.T) {
	doc, _ := Load(bytes.NewBufferString(str))

	doc.Set("key1", "new-value")
	newValue, _ := doc.Get("key1")
	expect(t, "修改已经存在的项的值", "new-value" == newValue)

	doc.Set("NOT-EXIST", "Setup-New-Item")
	newValue, _ = doc.Get("NOT-EXIST")
	expect(t, "修改不存在的项,默认是新增行为", "Setup-New-Item" == newValue)
}

func Test_Del(t *testing.T) {
	doc, _ := Load(bytes.NewBufferString(str))

	exist := doc.Del("NOT-EXIST")
	expect(t, "删除不存在的项,需要返回false", false == exist)

	exist = doc.Del("key1")
	expect(t, "删除已经存在的项,返回true", true == exist)
}

func Test_Comment_Uncomment(t *testing.T) {
	str := "key1=1\nkey 2 = 2"

	doc, _ := Load(bytes.NewBufferString(str))
	exist := doc.Comment("NOT-EXIST", "Some comment")
	expect(t, "对一个不存在的项执行注释操作,返回false", false == exist)

	doc.Comment("key1", "This is a \ncomment \nfor a")
	buf := bytes.NewBufferString("")
	err := doc.Save(buf)
	expect(t, "格式化成功", nil == err)
	exp1 := "#This is a \n#comment \n#for a\nkey1=1\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp1 == buf.String())

	doc.Comment("key 2", "")
	buf = bytes.NewBufferString("")
	err = doc.Save(buf)
	expect(t, "格式化成功", nil == err)
	exp2 := "#This is a \n#comment \n#for a\nkey1=1\n#\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp2 == buf.String())

	exist = doc.Uncomment("key1")
	expect(t, "对已经存在的key进行注释,返回true", true == exist)
	buf = bytes.NewBufferString("")
	err = doc.Save(buf)
	expect(t, "格式化成功", nil == err)
	exp3 := "key1=1\n#\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp3 == buf.String())

	exist = doc.Uncomment("key 2")
	expect(t, "对已经存在的key进行注释,返回true", true == exist)
	expect(t, "对已经存在的项进行注释", exp3 == buf.String())

	buf = bytes.NewBufferString("")
	err = doc.Save(buf)
	expect(t, "格式化成功", nil == err)
	exp4 := "key1=1\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp4 == buf.String())

	exist = doc.Uncomment("NOT-EXIST")
	expect(t, "对不已经存在的key进行注释,返回false", false == exist)
}

func Test_Accept(t *testing.T) {
	str := `
	#comment1
	key1=1
	#comment2
	key 2 : 2
	key3 : 2
	#coment3
	`

	count := 0

	doc, _ := Load(bytes.NewBufferString(str))
	doc.Accept(func(typo byte, value string, key string) bool {
		count++
		return typo != ':'
	})

	expect(t, "Accept提前中断", 5 == count)
}

func Test_Foreach(t *testing.T) {
	str := `
	#comment1
	key1=1
	#comment2
	key 2 : 2
	key3 : 2
	#coment3
	`

	count := 0

	doc, _ := Load(bytes.NewBufferString(str))
	doc.Foreach(func(value string, key string) bool {
		count++
		return "key 2" != key
	})

	expect(t, "Foreach提前中断", 2 == count)
}

func Test_IntDefault(t *testing.T) {
	str := `
	key1 : 1
	key2 : 123456789012345678901234567890
	key3 : abc
	key4 : 444 sina
	key5 : 555 007
	`
	doc, _ := Load(bytes.NewBufferString(str))

	v := doc.Int64Or("key1", 111)
	expect(t, "属性key已经存在时,返回存在的值", 1 == v)

	v = doc.Int64Or("NOT-EXIST", 111)
	expect(t, "属性key已经不存在时,返回缺省值", 111 == v)

	v = doc.Int64Or("key2", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)

	v = doc.Int64Or("key3", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)

	v = doc.Int64Or("key4", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)

	v = doc.Int64Or("key5", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)
}

func Test_UintDefault(t *testing.T) {
	str := `
	key1 : 1
	key2 : 123456789012345678901234567890
	key3 : abc
	key4 : 444 sina
	key5 : 555 007
	`
	doc, _ := Load(bytes.NewBufferString(str))

	v := doc.Uint64Or("key1", 111)
	expect(t, "属性key已经存在时,返回存在的值", 1 == v)

	v = doc.Uint64Or("NOT-EXIST", 111)
	expect(t, "属性key已经不存在时,返回缺省值", 111 == v)

	v = doc.Uint64Or("key2", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)

	v = doc.Uint64Or("key3", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)

	v = doc.Uint64Or("key4", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)

	v = doc.Uint64Or("key5", 111)
	expect(t, "属性key转换失败时,返回缺省值", 111 == v)
}

func Test_FloatDefault(t *testing.T) {
	str := `
	key1 : 1
	key2 : 123456789.012345678901234567890
	key3 : abc
	key4 : 123456789.a
	key5 : 123456789.
	`
	doc, _ := Load(bytes.NewBufferString(str))

	v := doc.Float64Or("key1", 111.0)
	expect(t, "属性key已经存在时,返回存在的值", 1.0 == v)

	v = doc.Float64Or("key2", 111.0)
	expect(t, "属性key已经存在时,返回存在的值", (v > 123456789.0) && (v < 123456789.1))

	v = doc.Float64Or("key5", 111.0)
	expect(t, "属性key已经存在时,返回存在的值", 123456789. == v)

	v = doc.Float64Or("key3", 111.0)
	expect(t, "转换失败时,返回def", 111.0 == v)

	v = doc.Float64Or("key4", 111.0)
	expect(t, "转换失败时,返回def", 111.0 == v)

	v = doc.Float64Or("NOT-EXIST", 111.0)
	expect(t, "属性不存在时,返回def", 111.0 == v)
}

func Test_BoolDefault(t *testing.T) {
	str := `
	key0 : 1
	key1 : T
	key2 : t
	key3 : true
	key4 : TRUE
	key5 : True
	key6 : 0
	key7 : F
	key8 : f
	key9 : false
	key10 : FALSE
	key11 : False
	key12 : Sina
	key13 : fALSE
	`

	doc, _ := Load(bytes.NewBufferString(str))

	for i := 0; i <= 5; i++ {
		v := doc.BoolOr(fmt.Sprintf("key%d", i), false)
		expect(t, "BoolDefault基本场景", v == true)
	}

	for i := 6; i <= 11; i++ {
		v := doc.BoolOr(fmt.Sprintf("key%d", i), true)
		expect(t, "BoolDefault基本场景", v == false)
	}

	v := doc.BoolOr("NOT-EXIST", true)
	expect(t, "获取不存在的项,返回def", v == true)

	v = doc.BoolOr("key12", false)
	expect(t, "无法转换的,返回def", v == false)

	v = doc.BoolOr("key13", true)
	expect(t, "无法转换的,返回def", v == true)
}

func Test_ObjectDefault(t *testing.T) {
	str := `
	key1 = 1
	key2 = man
	key3 = women
	key4 = 0
	key5
	`

	//	映射函数
	mapping := func(k string, v string) (interface{}, error) {
		if "0" == v {
			return 0, nil
		}

		if "1" == v {
			return 1, nil
		}

		if "man" == v {
			return 1, nil
		}

		if "women" == v {
			return 0, nil
		}

		if "" == v {
			return 0, errors.New("INVALID")
		}

		return -1, nil
	}

	doc, _ := Load(bytes.NewBufferString(str))

	expect(t, "ObjectOr:属性key已经存在时,返回存在的值1", 1 == doc.ObjectOr("key1", 123, mapping).(int))
	expect(t, "ObjectOr:属性key已经存在时,返回存在的值2", 1 == doc.ObjectOr("key2", 123, mapping).(int))
	expect(t, "ObjectOr:属性key已经存在时,返回存在的值3", 0 == doc.ObjectOr("key3", 123, mapping).(int))
	expect(t, "ObjectOr:属性key已经存在时,返回存在的值4", 0 == doc.ObjectOr("key4", 123, mapping).(int))
	expect(t, "ObjectOr:属性key转换失败时,返回def值5", 123 == doc.ObjectOr("key5", 123, mapping).(int))
	expect(t, "ObjectOr:属性key不存在时,返回def值5", 123 == doc.ObjectOr("NOT-EXIST", 123, mapping).(int))
}

func Test_Object(t *testing.T) {
	str := `
	key1 = 1
	key2 = man
	key3 = women
	key4 = 0
	key5
	`

	//	映射函数
	mapping := func(k string, v string) (interface{}, error) {
		if "0" == v {
			return 0, nil
		}

		if "1" == v {
			return 1, nil
		}

		if "man" == v {
			return 1, nil
		}

		if "women" == v {
			return 0, nil
		}

		if "" == v {
			return 0, errors.New("INVALID")
		}

		return -1, nil
	}

	doc, _ := Load(bytes.NewBufferString(str))

	expect(t, "ObjectOr:属性key已经存在时,返回存在的值1", 1 == doc.Object("key1", mapping).(int))

	//	nil不是万能类型,需要再想办法
	//expect(t, "ObjectOr:属性key不存在时,返回nil值1", 0 == doc.Object("NOT-EXIST", mapping).(int))
}

func Test_Int_String_Uint_Float_Bool(t *testing.T) {
	str := `
	key0 : -1
	key1 : timo
	key2 : 1234
	key3 : 12.5
	key4 : false
	`

	doc, _ := Load(bytes.NewBufferString(str))

	expect(t, "Int64", -1 == doc.Int64("key0"))
	expect(t, "String", "-1" == doc.String("key0"))
	expect(t, "String", "timo" == doc.String("key1"))
	expect(t, "String", "1234" == doc.String("key2"))
	expect(t, "String", "12.5" == doc.String("key3"))
	expect(t, "String", "false" == doc.String("key4"))
	expect(t, "Int64", 1234 == doc.Uint64("key2"))
	expect(t, "Int64", 12.5 == doc.Float64("key3"))
	expect(t, "Int64", false == doc.Bool("key4"))

}
