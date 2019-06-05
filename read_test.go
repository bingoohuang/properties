package properties

import (
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
	doc, _ := LoadString(str)

	v1 := doc.IntOr("key1", 111)
	expect(t, "属性key已经存在时,返回存在的值", v1 == 1)

	v1 = doc.IntOr("NOT-EXIST", 111)
	expect(t, "属性key已经不存在时,返回缺省值", v1 == 111)

	v := doc.Int64Or("key1", 111)
	expect(t, "属性key已经存在时,返回存在的值", v == 1)

	v = doc.Int64Or("NOT-EXIST", 111)
	expect(t, "属性key已经不存在时,返回缺省值", v == 111)

	v = doc.Int64Or("key2", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)

	v = doc.Int64Or("key3", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)

	v = doc.Int64Or("key4", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)

	v = doc.Int64Or("key5", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)
}

func Test_UintDefault(t *testing.T) {
	str := `
	key1 : 1
	key2 : 123456789012345678901234567890
	key3 : abc
	key4 : 444 sina
	key5 : 555 007
	`
	doc, _ := LoadString(str)

	v := doc.Uint64Or("key1", 111)
	expect(t, "属性key已经存在时,返回存在的值", v == 1)

	v = doc.Uint64Or("NOT-EXIST", 111)
	expect(t, "属性key已经不存在时,返回缺省值", v == 111)

	v = doc.Uint64Or("key2", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)

	v = doc.Uint64Or("key3", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)

	v = doc.Uint64Or("key4", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)

	v = doc.Uint64Or("key5", 111)
	expect(t, "属性key转换失败时,返回缺省值", v == 111)
}

func Test_FloatDefault(t *testing.T) {
	str := `
	key1 : 1
	key2 : 123456789.012345678901234567890
	key3 : abc
	key4 : 123456789.a
	key5 : 123456789.
	`
	doc, _ := LoadString(str)

	v := doc.Float64Or("key1", 111.0)
	expect(t, "属性key已经存在时,返回存在的值", v == 1.0)

	v = doc.Float64Or("key2", 111.0)
	expect(t, "属性key已经存在时,返回存在的值", v > 123456789.0 && v < 123456789.1)

	v = doc.Float64Or("key5", 111.0)
	expect(t, "属性key已经存在时,返回存在的值", v == 123456789.)

	v = doc.Float64Or("key3", 111.0)
	expect(t, "转换失败时,返回def", v == 111.0)

	v = doc.Float64Or("key4", 111.0)
	expect(t, "转换失败时,返回def", v == 111.0)

	v = doc.Float64Or("NOT-EXIST", 111.0)
	expect(t, "属性不存在时,返回def", v == 111.0)
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

	doc, _ := LoadString(str)

	for i := 0; i <= 5; i++ {
		v := doc.BoolOr(fmt.Sprintf("key%d", i), false)
		expect(t, "BoolDefault基本场景", v)
	}

	for i := 6; i <= 11; i++ {
		v := doc.BoolOr(fmt.Sprintf("key%d", i), true)
		expect(t, "BoolDefault基本场景", !v)
	}

	v := doc.BoolOr("NOT-EXIST", true)
	expect(t, "获取不存在的项,返回def", v)

	v = doc.BoolOr("key12", false)
	expect(t, "无法转换的,返回def", !v)

	v = doc.BoolOr("key13", true)
	expect(t, "无法转换的,返回def", v)
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
		switch v {
		case "0":
			return 0, nil
		case "1":
			return 1, nil
		case "man":
			return 1, nil
		case "women":
			return 0, nil
		case "":
			return 0, errors.New("INVALID")
		default:
			return -1, nil
		}
	}

	doc, _ := LoadString(str)

	expect(t, "ObjectOr:属性key已经存在时,返回存在的值1", doc.ObjectOr("key1", 123, mapping).(int) == 1)
	expect(t, "ObjectOr:属性key已经存在时,返回存在的值2", doc.ObjectOr("key2", 123, mapping).(int) == 1)
	expect(t, "ObjectOr:属性key已经存在时,返回存在的值3", doc.ObjectOr("key3", 123, mapping).(int) == 0)
	expect(t, "ObjectOr:属性key已经存在时,返回存在的值4", doc.ObjectOr("key4", 123, mapping).(int) == 0)
	expect(t, "ObjectOr:属性key转换失败时,返回def值5", doc.ObjectOr("key5", 123, mapping).(int) == 123)
	expect(t, "ObjectOr:属性key不存在时,返回def值5", doc.ObjectOr("NOT-EXIST", 123, mapping).(int) == 123)
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
	mapping := func(k, v string) (interface{}, error) {
		switch v {
		case "0":
			return 0, nil
		case "1":
			return 1, nil
		case "man":
			return 1, nil
		case "women":
			return 0, nil
		case "":
			return 0, errors.New("INVALID")
		default:
			return -1, nil
		}
	}

	doc, _ := LoadString(str)

	expect(t, "ObjectOr:属性key已经存在时,返回存在的值1", doc.Object("key1", mapping).(int) == 1)

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

	doc, _ := LoadString(str)

	expect(t, "Int64", doc.Int64("key0") == -1)
	expect(t, "Int", doc.Int("key0") == -1)
	expect(t, "Str", doc.Str("key0") == "-1")
	expect(t, "Str", doc.Str("key1") == "timo")
	expect(t, "Str", doc.Str("key2") == "1234")
	expect(t, "Str", doc.Str("key3") == "12.5")
	expect(t, "Str", doc.Str("key4") == "false")
	expect(t, "Int64", doc.Uint64("key2") == 1234)
	expect(t, "Int64", doc.Float64("key3") == 12.5)
	expect(t, "Int64", !doc.Bool("key4"))
}
