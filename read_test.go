package properties

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

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
	expect(t, "Int64", !doc.Bool("key4"))

}
